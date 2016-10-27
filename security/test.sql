CREATE TABLE qyx_usersource(
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  app_id VARCHAR(50) DEFAULT '' NOT NULL COMMENT 'APPID',
  open_id VARCHAR(50) DEFAULT '' NOT NULL COMMENT '用户ID',
  source_id VARCHAR(50) DEFAULT '' NOT NULL COMMENT '资源ID',
  `action` VARCHAR(50) DEFAULT '' NOT NULL COMMENT '行为',
  create_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  update_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间戳',
  KEY open_id (open_id),
  KEY app_id (app_id),
  KEY actions (action),
  KEY app_action_open_id (open_id,action,app_id)
)