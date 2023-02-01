create table if not exists `user` (
    `id` int primary key auto_increment comment 'user ID',
    `username` char(32) not null unique comment '注册用户名, 最长32个字符',
    `password` char(64) not null comment '密码, 加盐哈希后的',
    `follow_count` int default 0 comment '关注总数',
    `follower_count` int default 0 comment '粉丝总数'
) comment '用户' engine = innodb default charset = utf8;

create table if not exists `video` (
    `id` int primary key auto_increment comment 'video ID',
    `author_id` int not null comment '视频作者, vFK(user.id)',
    `title` varchar(255) not null comment '视频标题',
    `data` varchar(255) not null comment '视频数据, 文件系统路径',
    `cover` blob not null comment '视频封面图片',
    `favorite_count` int default 0 comment '点赞总数, 定时从 `favorite` 中异步统计更新',
    `comment_count` int default 0 comment '评论总数, 定时从 `comment` 中异步统计更新',
    index(`author_id`) comment '发布列表, 列出用户所有投稿过的视频'
) comment '视频' engine=innodb default charset=utf8;


create table if not exists `favorite` (
    `id` int primary key auto_increment comment 'favorite ID',
    `user_id` int not null comment '点赞的用户, vFK(user.id)',
    `video_id` int not null comment '被点赞的视频, vFK(video.id)',
    `time` timestamp default CURRENT_TIMESTAMP comment '点赞时间',
    index(`user_id`) comment '喜欢列表, 用户的所有点赞视频'
) comment '用户-点赞-视频' engine=innodb default charset=utf8;

create table if not exists `comment` (
    `id` int primary key auto_increment comment 'comment ID',
    `user_id` int not null comment '评论的用户, vFK(user.id)',
    `video_id` int not null comment '被评论的视频, vFK(video.id)',
    `time` timestamp default CURRENT_TIMESTAMP comment '评论时间',
    `comment_text` tinytext not null comment '评论内容',
    index(`video_id`) comment '评论列表, 查看视频的所有评论'
) comment '用户-评论-视频' engine=innodb default charset=utf8;

create table if not exists `relation` (
    `id` int primary key auto_increment comment 'relation ID',
    `user_id` int not null comment 'vFK(user.id)',
    `followed_user_id` int not null comment '被关注的用户, vFK(user.id)',
    `time` timestamp default CURRENT_TIMESTAMP comment '关注时间',
    index(`user_id`),
    index(`followed_user_id`)
) comment '用户-关注-用户' engine=innodb default charset=utf8;

create table if not exists `message` (
    `id` int primary key auto_increment comment 'message ID',
    `send_user_id` int not null comment '发送消息的用户, vFK(user.id)',
    `receive_user_id` int not null comment '接收消息的用户, vFK(user.id)',
    `time` timestamp default CURRENT_TIMESTAMP comment '消息发送时间',
    `content` tinytext not null comment '消息内容',
    index(`send_user_id`),
    index(`receive_user_id`)
) comment '用户-发消息-用户' engine=innodb default charset=utf8;