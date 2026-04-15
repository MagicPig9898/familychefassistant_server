# -*- coding: utf-8 -*-
"""
百度百科蔬菜食材爬虫
- 从百度百科爬取 200 个常见蔬菜类食材的详细信息
- 输出为 Markdown 风格的 txt 文件，方便后续 RAG 切分
"""

import os
import sys
import time
import random
import urllib.parse
import urllib.request
import urllib.error
import json
import re
from html.parser import HTMLParser

# ============================================================
# 配置
# ============================================================
RAW_DIR = os.path.join(os.path.dirname(os.path.abspath(__file__)), "raw")
BASE_URL = "https://baike.baidu.com/item/"

HEADERS = {
    "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
    "Accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8",
    "Accept-Language": "zh-CN,zh;q=0.9,en;q=0.8",
    "Accept-Encoding": "identity",
    "Connection": "keep-alive",
    "Cache-Control": "max-age=0",
}

# 随机延时范围（秒）
DELAY_MIN = 2.0
DELAY_MAX = 5.0


# ============================================================
# 轻量 HTML -> 文本 解析器（不依赖第三方库）
# ============================================================
class BaikeParser(HTMLParser):
    """解析百度百科 HTML，提取结构化文本内容"""

    def __init__(self):
        super().__init__()
        self.result = []
        self.current_section = []
        self.in_tag = None
        self.tag_stack = []
        self.skip_tags = {"script", "style", "noscript", "svg", "path"}
        self.skip_depth = 0

        # 百度百科相关 CSS 类
        self.in_summary = False
        self.in_basic_info = False
        self.in_content = False
        self.in_para = False

        # 基本信息栏
        self.basic_info_pairs = []
        self.current_dt = ""
        self.current_dd = ""
        self.in_dt = False
        self.in_dd = False

        # 标题
        self.in_heading = False
        self.heading_level = 0
        self.heading_text = ""

        # 段落文本
        self.para_text = ""

    def handle_starttag(self, tag, attrs):
        attrs_dict = dict(attrs)
        class_name = attrs_dict.get("class", "")

        # 跳过不需要的标签
        if tag in self.skip_tags:
            self.skip_depth += 1
            return

        if self.skip_depth > 0:
            return

        self.tag_stack.append(tag)

        # 摘要区域
        if "lemma-summary" in class_name or "lemmaWgt-lemmaSummary" in class_name:
            self.in_summary = True
            self.para_text = ""

        # 基本信息栏
        if "basic-info" in class_name or "basicInfo-block" in class_name:
            self.in_basic_info = True

        # dt / dd
        if tag == "dt" and self.in_basic_info:
            self.in_dt = True
            self.current_dt = ""
        if tag == "dd" and self.in_basic_info:
            self.in_dd = True
            self.current_dd = ""

        # 正文内容区
        if "para-title" in class_name or "paraTitle" in class_name:
            self.in_heading = True
            self.heading_text = ""
            # 判断层级
            if tag in ("h2",):
                self.heading_level = 2
            elif tag in ("h3",):
                self.heading_level = 3
            else:
                # 通过 class 判断
                if "level-2" in class_name:
                    self.heading_level = 2
                elif "level-3" in class_name:
                    self.heading_level = 3
                else:
                    self.heading_level = 2

        if tag in ("h2", "h3") and not self.in_heading:
            self.in_heading = True
            self.heading_text = ""
            self.heading_level = 2 if tag == "h2" else 3

        # 段落
        if "para" in class_name.split() or (tag == "div" and "para" in class_name):
            self.in_para = True
            self.para_text = ""

    def handle_endtag(self, tag):
        if tag in self.skip_tags:
            self.skip_depth -= 1
            return

        if self.skip_depth > 0:
            return

        if self.tag_stack and self.tag_stack[-1] == tag:
            self.tag_stack.pop()

        # 摘要结束
        if self.in_summary and tag in ("div",):
            text = self.para_text.strip()
            if text and len(text) > 10:
                self.result.append(("summary", text))
                self.in_summary = False
                self.para_text = ""

        # dt/dd 结束
        if tag == "dt" and self.in_dt:
            self.in_dt = False
        if tag == "dd" and self.in_dd:
            self.in_dd = False
            dt = self.current_dt.strip().rstrip("：:").strip()
            dd = self.current_dd.strip()
            if dt and dd:
                self.basic_info_pairs.append((dt, dd))

        # 基本信息栏结束
        if self.in_basic_info and tag in ("div",) and not self.in_dt and not self.in_dd:
            if self.basic_info_pairs:
                self.result.append(("basic_info", list(self.basic_info_pairs)))
                self.in_basic_info = False

        # 标题结束
        if self.in_heading and tag in ("h2", "h3", "div", "span"):
            text = self.heading_text.strip()
            if text and len(text) < 50:
                self.result.append(("heading", self.heading_level, text))
            self.in_heading = False
            self.heading_text = ""

        # 段落结束
        if self.in_para and tag == "div":
            text = self.para_text.strip()
            if text and len(text) > 5:
                self.result.append(("para", text))
            self.in_para = False
            self.para_text = ""

    def handle_data(self, data):
        if self.skip_depth > 0:
            return

        text = data.strip()
        if not text:
            return

        if self.in_summary:
            self.para_text += data
        if self.in_dt:
            self.current_dt += data
        if self.in_dd:
            self.current_dd += data
        if self.in_heading:
            self.heading_text += data
        if self.in_para:
            self.para_text += data


# ============================================================
# 备用解析：正则表达式方式
# ============================================================
def parse_with_regex(html: str, name: str) -> str:
    """使用正则表达式从百度百科 HTML 中提取内容"""
    lines = [f"# {name}\n"]

    # 去掉 script / style
    html_clean = re.sub(r"<script[^>]*>.*?</script>", "", html, flags=re.DOTALL | re.IGNORECASE)
    html_clean = re.sub(r"<style[^>]*>.*?</style>", "", html_clean, flags=re.DOTALL | re.IGNORECASE)

    def strip_tags(s):
        return re.sub(r"<[^>]+>", "", s).strip()

    # 1. 摘要
    summary_match = re.search(
        r'<div class="[^"]*lemma-summary[^"]*"[^>]*>(.*?)</div>\s*(?:<div class="(?!.*lemma-summary))',
        html_clean, re.DOTALL | re.IGNORECASE
    )
    if not summary_match:
        summary_match = re.search(
            r'<div class="[^"]*lemmaSummary[^"]*"[^>]*>(.*?)</div>',
            html_clean, re.DOTALL | re.IGNORECASE
        )
    if not summary_match:
        # 尝试 meta description
        meta_match = re.search(r'<meta\s+name="description"\s+content="([^"]+)"', html_clean, re.IGNORECASE)
        if meta_match:
            lines.append("## 摘要\n")
            lines.append(meta_match.group(1).strip() + "\n")
    else:
        text = strip_tags(summary_match.group(1))
        if text:
            lines.append("## 摘要\n")
            lines.append(text + "\n")

    # 2. 基本信息
    basic_match = re.search(
        r'<div class="[^"]*basic-info[^"]*"[^>]*>(.*?)</div>\s*</div>',
        html_clean, re.DOTALL | re.IGNORECASE
    )
    if basic_match:
        basic_html = basic_match.group(1)
        dts = re.findall(r"<dt[^>]*>(.*?)</dt>", basic_html, re.DOTALL)
        dds = re.findall(r"<dd[^>]*>(.*?)</dd>", basic_html, re.DOTALL)
        if dts and dds:
            lines.append("## 基本信息\n")
            for dt, dd in zip(dts, dds):
                dt_text = strip_tags(dt).rstrip("：:").strip()
                dd_text = strip_tags(dd).strip()
                if dt_text and dd_text:
                    lines.append(f"- {dt_text}: {dd_text}")
            lines.append("")

    # 3. 正文章节 - 提取所有标题和段落
    # 查找 h2/h3 标题
    heading_pattern = re.compile(
        r'<(h[23])[^>]*class="[^"]*"[^>]*>(.*?)</\1>|'
        r'<div[^>]*class="[^"]*para-title[^"]*"[^>]*>(.*?)</div>',
        re.DOTALL | re.IGNORECASE
    )
    para_pattern = re.compile(
        r'<div[^>]*class="[^"]*\bpara\b[^"]*"[^>]*>(.*?)</div>',
        re.DOTALL | re.IGNORECASE
    )

    # 按位置合并标题和段落
    elements = []
    for m in heading_pattern.finditer(html_clean):
        text = strip_tags(m.group(2) or m.group(3) or "")
        if text and len(text) < 50:
            tag = m.group(1) or "h2"
            level = 2 if "2" in tag else 3
            elements.append((m.start(), "heading", level, text))

    for m in para_pattern.finditer(html_clean):
        text = strip_tags(m.group(1))
        if text and len(text) > 5:
            elements.append((m.start(), "para", text))

    elements.sort(key=lambda x: x[0])

    for elem in elements:
        if elem[1] == "heading":
            prefix = "#" * (elem[2] + 1)  # h2 -> ###, h3 -> ####
            lines.append(f"\n{prefix} {elem[3]}\n")
        elif elem[1] == "para":
            lines.append(elem[2] + "\n")

    # 4. 如果正文太少，尝试提取所有可见文本
    if len(lines) < 5:
        # 从 meta 标签获取信息
        keywords_match = re.search(r'<meta\s+name="keywords"\s+content="([^"]+)"', html_clean, re.IGNORECASE)
        if keywords_match:
            lines.append(f"\n## 关键词\n")
            lines.append(keywords_match.group(1).strip() + "\n")

        desc_match = re.search(r'<meta\s+name="description"\s+content="([^"]+)"', html_clean, re.IGNORECASE)
        if desc_match and "## 摘要" not in "\n".join(lines):
            lines.append(f"\n## 描述\n")
            lines.append(desc_match.group(1).strip() + "\n")

        # 提取 JSON-LD 结构化数据
        jsonld_match = re.search(r'<script[^>]*type="application/ld\+json"[^>]*>(.*?)</script>', html, re.DOTALL)
        if jsonld_match:
            try:
                data = json.loads(jsonld_match.group(1))
                if isinstance(data, dict):
                    if "description" in data:
                        lines.append(f"\n## 详细描述\n")
                        lines.append(data["description"].strip() + "\n")
            except (json.JSONDecodeError, KeyError):
                pass

    return "\n".join(lines)


# ============================================================
# 核心函数
# ============================================================
def fetch_baike(name: str) -> str | None:
    """请求百度百科词条页面，返回 HTML"""
    encoded_name = urllib.parse.quote(name)
    url = BASE_URL + encoded_name

    req = urllib.request.Request(url, headers=HEADERS)
    try:
        with urllib.request.urlopen(req, timeout=15) as resp:
            # 处理编码
            content_type = resp.headers.get("Content-Type", "")
            if "charset=" in content_type:
                encoding = content_type.split("charset=")[-1].strip()
            else:
                encoding = "utf-8"
            html = resp.read().decode(encoding, errors="replace")
            return html
    except urllib.error.HTTPError as e:
        print(f"  HTTP 错误 {e.code}: {url}")
        return None
    except urllib.error.URLError as e:
        print(f"  URL 错误: {e.reason}")
        return None
    except Exception as e:
        print(f"  请求异常: {e}")
        return None


def parse_baike(html: str, name: str) -> str:
    """解析百度百科 HTML，提取结构化文本"""
    # 优先用 HTMLParser
    try:
        parser = BaikeParser()
        parser.feed(html)

        lines = [f"# {name}\n"]

        for item in parser.result:
            if item[0] == "summary":
                lines.append("## 摘要\n")
                lines.append(item[1] + "\n")
            elif item[0] == "basic_info":
                lines.append("## 基本信息\n")
                for dt, dd in item[1]:
                    lines.append(f"- {dt}: {dd}")
                lines.append("")
            elif item[0] == "heading":
                prefix = "#" * (item[1] + 1)
                lines.append(f"\n{prefix} {item[2]}\n")
            elif item[0] == "para":
                lines.append(item[1] + "\n")

        text = "\n".join(lines)

        # 如果 HTMLParser 结果太少，降级到正则
        if len(lines) < 5:
            text = parse_with_regex(html, name)

    except Exception:
        text = parse_with_regex(html, name)

    # 清理多余空行
    text = re.sub(r"\n{3,}", "\n\n", text)
    return text.strip()


def save_to_file(filename: str, content: str):
    """保存到 raw/ 目录"""
    filepath = os.path.join(RAW_DIR, f"{filename}.txt")
    with open(filepath, "w", encoding="utf-8") as f:
        f.write(content)


def main():
    # 导入蔬菜列表
    from vegetables_list import VEGETABLES

    os.makedirs(RAW_DIR, exist_ok=True)

    total = len(VEGETABLES)
    success_count = 0
    skip_count = 0
    fail_count = 0
    fail_list = []

    print("=" * 60)
    print("  baidu baike vegetable crawler")
    print(f"  total: {total}")
    print(f"  output: {RAW_DIR}")
    print("=" * 60)
    print()

    for i, (cn_name, en_name) in enumerate(VEGETABLES):
        filepath = os.path.join(RAW_DIR, f"{en_name}.txt")

        # 断点续爬：已存在的文件跳过
        if os.path.exists(filepath):
            file_size = os.path.getsize(filepath)
            if file_size > 50:  # 内容大于 50 字节才算有效
                skip_count += 1
                print(f"[{i + 1}/{total}] skip: {cn_name} ({en_name}.txt, {file_size}B)")
                continue

        print(f"[{i + 1}/{total}] crawling: {cn_name} ...", end=" ", flush=True)

        html = fetch_baike(cn_name)
        if html:
            content = parse_baike(html, cn_name)
            content_len = len(content)

            if content_len > 50:
                save_to_file(en_name, content)
                success_count += 1
                print(f"[OK] saved: raw/{en_name}.txt ({content_len} chars)")
            else:
                # 内容太少，可能词条不存在或解析失败
                # 依然保存，但标记
                save_to_file(en_name, content)
                fail_count += 1
                fail_list.append(cn_name)
                print(f"[WARN] too short: raw/{en_name}.txt ({content_len} chars)")
        else:
            fail_count += 1
            fail_list.append(cn_name)
            print(f"[FAIL] fetch failed")

        # 随机延时，避免被反爬
        delay = random.uniform(DELAY_MIN, DELAY_MAX)
        time.sleep(delay)

    # 统计
    print()
    print("=" * 60)
    print("  Done!")
    print(f"  Success: {success_count}")
    print(f"  Skipped: {skip_count}")
    print(f"  Failed:  {fail_count}")
    if fail_list:
        print(f"  Failed items: {', '.join(fail_list)}")
    print("=" * 60)


if __name__ == "__main__":
    main()
