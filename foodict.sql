-- ========================================
-- 数据库：food_dict
-- ========================================
CREATE DATABASE IF NOT EXISTS food_dict
DEFAULT CHARACTER SET utf8mb4
COLLATE utf8mb4_unicode_ci;

USE food_dict;

-- ========================================
-- 表1：食材主表
-- ========================================
CREATE TABLE ingredients (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '主键ID',
    name VARCHAR(50) COMMENT '食材名称',
    alias_name VARCHAR(50) COMMENT '别名：str数组的json字符串',
    category_id BIGINT COMMENT '分类',
    flavor VARCHAR(100) COMMENT '口味，str数组的json字符串',
    properties VARCHAR(100) COMMENT '属性，str数组的json字符串',
    season VARCHAR(100) COMMENT '季节，str数组的json字符串',
    created_at BIGINT COMMENT '创建时间',
    updated_at BIGINT COMMENT '更新时间'
) COMMENT='食材主表';

-- ========================================
-- 表2：食材分类表
-- ========================================
CREATE TABLE ingredient_category (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '主键ID',
    category VARCHAR(50) NOT NULL COMMENT '分类'
) COMMENT='食材分类表';


-- ========================================
-- 表3：营养信息表
-- ========================================
CREATE TABLE ingredient_nutrition (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '主键ID',
    ingredient_id BIGINT COMMENT '食材ID',
    summary TEXT COMMENT '营养提要, 如富含维生素C...使用语义完整的一句话, 内容多可以分多行记录',
    INDEX idx_ingredient_id (ingredient_id)
) COMMENT='食材营养信息表';

-- ========================================
-- 表4：适宜/不适人群表
-- ========================================
CREATE TABLE ingredient_suitable (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '主键ID',
    ingredient_id BIGINT COMMENT '食材ID',
    summary TEXT COMMENT '提要，如：减肥人群，因为...使用语义完整的一句话, 内容多可以分多行记录',
    type TINYINT COMMENT '类型: 1=适宜,0=不适宜',
    INDEX idx_ingredient_id (ingredient_id)
) COMMENT='适宜/不适人群表';

-- ========================================
-- 表5：食材搭配表
-- ========================================
CREATE TABLE ingredient_pairing (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '主键ID',
    ingredient_id BIGINT COMMENT '当前食材ID',
    summary TEXT COMMENT '搭配提要,如:搭配豆腐，因为...使用语义完整的一句话, 内容多可以分多行记录',
    type TINYINT COMMENT '类型: 1=推荐搭配,0=不推荐搭配',
    INDEX idx_ingredient_id (ingredient_id)
) COMMENT='食材搭配关系表';


-- ========================================
-- 表6：为RAG预留的扩展
-- ========================================
CREATE TABLE ingredient_chunks (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    ingredient_id BIGINT COMMENT '食材ID',
    content TEXT COMMENT '切分后的文本块',
    embedding JSON COMMENT '向量（如果数据库支持）',
    type TINYINT COMMENT '类型:0=nutrition/1=pairing/2=suitable',
    INDEX idx_ingredient_id (ingredient_id)
) COMMENT='食材RAG文本chunk表';