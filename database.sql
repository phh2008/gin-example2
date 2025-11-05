-- ----------------------------
-- Table structure for sys_permission
-- ----------------------------
DROP TABLE IF EXISTS `sys_permission`;
CREATE TABLE `sys_permission`
(
    `id`         bigint       NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `perm_name`  varchar(50)  NOT NULL DEFAULT '' COMMENT '权限名称',
    `url`        varchar(200) NOT NULL DEFAULT '' COMMENT 'URL路径',
    `action`     varchar(50)  NOT NULL DEFAULT '' COMMENT '权限动作：比如get、post、delete等',
    `perm_type`  tinyint      NOT NULL DEFAULT 1 COMMENT '权限类型：1-菜单、2-按钮',
    `parent_id`  bigint       NOT NULL DEFAULT 0 COMMENT '父级ID：资源层级关系',
    `created_at` datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
    `created_by` varchar(50) NULL DEFAULT '' COMMENT '创建人',
    `updated_by` varchar(50) NULL DEFAULT '' COMMENT '更新人',
    `deleted`    tinyint      NOT NULL DEFAULT 1 COMMENT '是否删除：1-否，2-是',
    PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB COMMENT = '权限表';


-- ----------------------------
-- Table structure for sys_role
-- ----------------------------
DROP TABLE IF EXISTS `sys_role`;
CREATE TABLE `sys_role`
(
    `id`         bigint      NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `role_code`  varchar(50) NOT NULL DEFAULT '' COMMENT '角色编号',
    `role_name`  varchar(50) NOT NULL DEFAULT '' COMMENT '角色名称',
    `created_at` datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
    `created_by` varchar(50) NULL DEFAULT '' COMMENT '创建人',
    `updated_by` varchar(50) NULL DEFAULT '' COMMENT '更新人',
    `deleted`    tinyint     NOT NULL DEFAULT 1 COMMENT '是否删除：1-否，2-是',
    PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB COMMENT = '角色表';

-- ----------------------------
-- Table structure for sys_role_permission
-- ----------------------------
DROP TABLE IF EXISTS `sys_role_permission`;
CREATE TABLE `sys_role_permission`
(
    `id`      bigint NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `role_id` bigint NOT NULL DEFAULT 0 COMMENT '角色编号',
    `perm_id` bigint NOT NULL DEFAULT 0 COMMENT '权限ID',
    PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB COMMENT = '角色权限表';


-- ----------------------------
-- Table structure for sys_user
-- ----------------------------
DROP TABLE IF EXISTS `sys_user`;
CREATE TABLE `sys_user`
(
    `id`         bigint       NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `real_name`  varchar(50)  NOT NULL DEFAULT '' COMMENT '姓名',
    `user_name`  varchar(50)  NOT NULL DEFAULT '' COMMENT '用户名',
    `email`      varchar(50)  NOT NULL DEFAULT '' COMMENT '邮箱',
    `password`   varchar(200) NOT NULL DEFAULT '' COMMENT '密码',
    `status`     tinyint      NOT NULL DEFAULT 1 COMMENT '状态: 1-启用，2-禁用',
    `role_code`  varchar(50)  NOT NULL DEFAULT '' COMMENT '角色编号',
    `created_at` datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
    `created_by` varchar(50) NULL DEFAULT '' COMMENT '创建人',
    `updated_by` varchar(50) NULL DEFAULT '' COMMENT '更新人',
    `deleted`    tinyint      NOT NULL DEFAULT 1 COMMENT '是否删除：1-否，2-是',
    PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB COMMENT = '用户表';

CREATE TABLE `sys_user_role`
(
    `id`      BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `user_id` BIGINT NOT NULL DEFAULT '0' COMMENT '用户ID',
    `role_id` BIGINT NOT NULL DEFAULT '0' COMMENT '角色ID',
    PRIMARY KEY (`id`),
    INDEX     `idx_user_id_role_id` (`user_id`, `role_id`)
) COMMENT='用户角色表' ENGINE=InnoDB;
