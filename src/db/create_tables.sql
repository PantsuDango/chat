CREATE DATABASE `chat_db`;
USE `chat_db`;

-- 聊天消息记录表
CREATE TABLE `chat_message` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `ip` varchar(32) NOT NULL COMMENT 'IP地址',
    `message` text NOT NULL COMMENT '聊天消息',
    `message_type` enum('First', 'Option', 'Manual', 'Keyword', 'Customer') NOT NULL COMMENT '消息类型: 首次/质保/人工/关键词/客户',
    `createtime` datetime NOT NULL COMMENT '创建时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY (`ip`),
    KEY (`ip`, `createtime`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COMMENT='聊天消息记录表';

-- IP备注表
CREATE TABLE `ip_content_map` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `ip` varchar(32) NOT NULL COMMENT 'IP地址',
    `content` text NOT NULL COMMENT '备注信息',
    `createtime` datetime NOT NULL COMMENT '创建时间',
    `lastupdate` datetime NOT NULL COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY (`ip`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COMMENT='IP备注表';

-- 首次回复设置表
CREATE TABLE `first_reply` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `message` text DEFAULT NULL COMMENT '回复消息',
    `option_switch` tinyint(4) DEFAULT 0 COMMENT '选项开关: 0-关, 1-开',
    `createtime` datetime NOT NULL COMMENT '创建时间',
    `lastupdate` datetime NOT NULL COMMENT '更新时间',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COMMENT='首次回复设置表';
INSERT INTO `first_reply` VALUES (1, NULL, 0, NOW(), NOW());

-- 首次回复选项消息表
CREATE TABLE `first_reply_option_message` (
   `id` int(11) NOT NULL AUTO_INCREMENT,
   `option` varchar(128) NOT NULL COMMENT '附带选项',
   `content` text DEFAULT NULL COMMENT '选项消息',
   `createtime` datetime NOT NULL COMMENT '创建时间',
   `lastupdate` datetime NOT NULL COMMENT '更新时间',
   PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COMMENT='首次回复选项消息表';
INSERT INTO `first_reply_option_message` VALUES (1, '质保', NULL, NOW(), NOW());

-- 关键词规则
CREATE TABLE `keyword_rule` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `rule_name` varchar(128) NOT NULL COMMENT '规则名称',
    `switch` tinyint(4) DEFAULT 0 COMMENT '规则开关: 0-关, 1-开',
    `content` text NOT NULL COMMENT '回复消息',
    `createtime` datetime NOT NULL COMMENT '创建时间',
    `lastupdate` datetime NOT NULL COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY (`rule_name`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COMMENT='关键词规则';

-- 规则与关键词映射表
CREATE TABLE `keyword_rule_map` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `rule_name` varchar(128) NOT NULL COMMENT '规则名称',
    `content` varchar(256) NOT NULL COMMENT '关键词',
    `createtime` datetime NOT NULL COMMENT '创建时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY (`rule_name`, `content`),
    KEY (`rule_name`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COMMENT='规则与关键词映射表';