create table if not exists `user` (
    `id` int primary key auto_increment comment 'user ID',
    `username` char(64) not null unique comment '注册用户名, 最长32个字符',
    `password` char(60) not null comment '密码, 加盐哈希后的',
    `follow_count` int default 0 comment '关注总数',
    `follower_count` int default 0 comment '粉丝总数'
    # index(`username`) # unique 约束会创建一个唯一索引
) comment '用户' engine = innodb default charset = utf8;

create table if not exists `video` (
    `id` int primary key auto_increment comment 'video ID',
    `author_id` int not null comment '视频作者, vFK(user.id)',
    `title` varchar(255) not null comment '视频标题',
    `data` varchar(255) not null comment '视频数据, 文件系统路径',
    `cover` varchar(255) not null comment '视频封面图片, 文件系统路径',
    `favorite_count` int default 0 comment '点赞总数, 定时从 `favorite` 中异步统计更新',
    `comment_count` int default 0 comment '评论总数, 定时从 `comment` 中异步统计更新',
    `time` timestamp default CURRENT_TIMESTAMP comment '投稿时间',
    index(`author_id`) comment '发布列表, 列出用户所有投稿过的视频',
    index(`time`) comment '视频流接口, 按投稿时间倒序的视频列表'
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


# TEST DATA

insert into `user` (`id`, `username`, `follow_count`, `follower_count`, `password`) values
    (1, 'TEST_aBadString', 3, 1, '$2a$10$kWnGFTqsXTHcETyIQxvxD.iMpqrE2oMtxAAWXQPbE0JN1wC4IDn0q'),
    (2, 'TEST_peadx', 2, 1, '$2a$10$kWnGFTqsXTHcETyIQxvxD.iMpqrE2oMtxAAWXQPbE0JN1wC4IDn0q'),
    (3, 'TEST_bin', 1, 2, '$2a$10$kWnGFTqsXTHcETyIQxvxD.iMpqrE2oMtxAAWXQPbE0JN1wC4IDn0q'),
    (4, 'TEST_song', 0, 2, '$2a$10$kWnGFTqsXTHcETyIQxvxD.iMpqrE2oMtxAAWXQPbE0JN1wC4IDn0q');

insert into `video` (`id`, `author_id`, `title`, `data`, `cover`, `favorite_count`) values
    (1, 1, 'TEST_aBadString_video_1', 'default.mp4', 'default.jpg', 4),
    (2, 1, 'TEST_aBadString_video_2', 'default.mp4', 'default.jpg', 1),
    (3, 1, 'TEST_aBadString_video_3', 'default.mp4', 'default.jpg', 1),
    (4, 2, 'TEST_peadx_video_1', 'default.mp4', 'default.jpg', 1),
    (5, 2, 'TEST_peadx_video_2', 'default.mp4', 'default.jpg', 1),
    (6, 2, 'TEST_peadx_video_3', 'default.mp4', 'default.jpg', 0);

insert into `relation` (`user_id`, `followed_user_id`) values
    (1, 2),
    (1, 3),
    (1, 4),
    (2, 1),
    (2, 3),
    (3, 4);

insert into `favorite` (user_id, video_id) values
    (1, 1),
    (1, 2),
    (1, 3),
    (1, 4),
    (1, 5),
    (2, 1),
    (3, 1),
    (4, 1);