-- phpMyAdmin SQL Dump
-- version 4.5.1
-- http://www.phpmyadmin.net
--
-- Host: 127.0.0.1
-- Generation Time: 2017-05-06 13:35:58
-- 服务器版本： 10.1.10-MariaDB
-- PHP Version: 7.0.4

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- Database: `webim`
--

-- --------------------------------------------------------

--
-- 表的结构 `contacts`
--

CREATE TABLE `contacts` (
  `id` int(11) NOT NULL,
  `master_uid` int(11) UNSIGNED NOT NULL,
  `group_id` int(10) UNSIGNED NOT NULL,
  `uid` int(11) UNSIGNED NOT NULL,
  `order_weight` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

--
-- 转存表中的数据 `contacts`
--

INSERT INTO `contacts` (`id`, `master_uid`, `group_id`, `uid`, `order_weight`) VALUES
(1, 4, 1, 1, 0),
(2, 4, 1, 2, 0),
(3, 4, 1, 7, 0),
(4, 7, 4, 4, 0),
(5, 1111, 2222, 33333, 0),
(6, 1, 0, 7, 0),
(7, 7, 4, 1, 0);

-- --------------------------------------------------------

--
-- 表的结构 `contact_group`
--

CREATE TABLE `contact_group` (
  `id` int(11) NOT NULL,
  `uid` int(10) UNSIGNED NOT NULL,
  `title` varchar(20) NOT NULL DEFAULT '',
  `order_weight` int(10) UNSIGNED NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

--
-- 转存表中的数据 `contact_group`
--

INSERT INTO `contact_group` (`id`, `uid`, `title`, `order_weight`) VALUES
(1, 4, '网红', 0),
(2, 4, '前端码屌', 0),
(3, 4, '我心中的女神', 0),
(4, 7, '默认', 0),
(5, 111, '实打实的', 0);

-- --------------------------------------------------------

--
-- 表的结构 `global_group`
--

CREATE TABLE `global_group` (
  `id` int(11) NOT NULL,
  `title` varchar(20) NOT NULL DEFAULT '',
  `channel_id` varchar(32) NOT NULL DEFAULT '',
  `pic` varchar(120) NOT NULL DEFAULT '',
  `remark` varchar(500) NOT NULL DEFAULT ''
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

--
-- 转存表中的数据 `global_group`
--

INSERT INTO `global_group` (`id`, `title`, `channel_id`, `pic`, `remark`) VALUES
(1, '前端群', 'channel_id_1', 'http://tp2.sinaimg.cn/2211874245/180/40050524279/0', ''),
(2, 'Fly社区官方群', 'channel_id_2', 'http://tp2.sinaimg.cn/5488749285/50/5719808192/1', ''),
(3, '前端群2', 'channel_id_3', 'http://tp2.sinaimg.cn/2211874245/180/40050524279/0', ''),
(4, '前端群4', 'channel_id_4', 'http://tp2.sinaimg.cn/2211874245/180/40050524279/0', ''),
(5, '前端群5', 'channel_id_5', 'http://tp2.sinaimg.cn/2211874245/180/40050524279/0', '');

-- --------------------------------------------------------

--
-- 表的结构 `req_friend`
--

CREATE TABLE `req_friend` (
  `id` int(10) UNSIGNED NOT NULL,
  `type` tinyint(1) UNSIGNED NOT NULL DEFAULT '1' COMMENT '1申请添加好友,2系统通知',
  `from_uid` int(10) UNSIGNED NOT NULL,
  `add_group` int(10) UNSIGNED NOT NULL COMMENT '添加到哪个分组',
  `readed` tinyint(1) UNSIGNED NOT NULL DEFAULT '0',
  `req_uid` int(10) UNSIGNED NOT NULL,
  `status` int(10) UNSIGNED NOT NULL COMMENT '1申请中,2同意,3拒绝',
  `remark` varchar(50) NOT NULL DEFAULT '',
  `time` int(10) UNSIGNED NOT NULL,
  `up_time` int(10) UNSIGNED NOT NULL DEFAULT '0'
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

--
-- 转存表中的数据 `req_friend`
--

INSERT INTO `req_friend` (`id`, `type`, `from_uid`, `add_group`, `readed`, `req_uid`, `status`, `remark`, `time`, `up_time`) VALUES
(1, 1, 1, 0, 0, 7, 2, '', 0, 0),
(2, 1, 2, 1, 0, 7, 2, '', 0, 0),
(3, 1, 1, 0, 0, 0, 1, '', 1493294893, 0),
(4, 1, 1, 0, 0, 0, 1, '', 1493294931, 0),
(5, 1, 1, 0, 0, 0, 1, '', 1493295001, 0),
(6, 1, 1, 0, 0, 0, 1, '', 1493295074, 0),
(7, 1, 1, 0, 0, 0, 1, '', 1493295118, 0),
(8, 1, 1, 0, 0, 0, 1, '', 1493295160, 0),
(9, 1, 1, 0, 0, 0, 1, '', 1493295166, 0),
(10, 1, 1, 0, 0, 0, 1, '', 1493347415, 1493347415),
(11, 1, 1, 0, 0, 0, 1, '', 1493717747, 1493717747),
(12, 1, 1, 0, 0, 0, 1, '', 1493717753, 1493717753);

-- --------------------------------------------------------

--
-- 表的结构 `user`
--

CREATE TABLE `user` (
  `id` int(11) NOT NULL,
  `user` varchar(20) NOT NULL DEFAULT '',
  `pwd` varchar(32) NOT NULL DEFAULT '',
  `sid` varchar(32) NOT NULL DEFAULT '',
  `nick` varchar(20) NOT NULL DEFAULT '',
  `age` int(10) UNSIGNED NOT NULL,
  `sign` varchar(200) NOT NULL DEFAULT '',
  `reg_time` int(10) UNSIGNED NOT NULL,
  `is_online` tinyint(1) UNSIGNED NOT NULL DEFAULT '0',
  `status` varchar(20) NOT NULL DEFAULT 'offline',
  `avatar` varchar(120) NOT NULL DEFAULT '',
  `token` varchar(128) NOT NULL DEFAULT ''
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

--
-- 转存表中的数据 `user`
--

INSERT INTO `user` (`id`, `user`, `pwd`, `sid`, `nick`, `age`, `sign`, `reg_time`, `is_online`, `status`, `avatar`, `token`) VALUES
(1, 'user', '123456', '', 'nick1', 0, '', 0, 0, 'offline', '/im/avatar/femalecodertocat.png', ''),
(2, 'user2', '123456', '', 'nick2', 0, '', 0, 0, 'offline', '/im/avatar/mountietocat.png', ''),
(4, 'weichaoduo', '121', '1950501661', '纸飞机', 32, '在深邃的编码世界，做一枚轻盈的纸飞机', 1492770152, 0, 'offline', '/im/avatar/octocat-de-los-muertos.jpg', '62980976412'),
(5, 'weichaoduo2', '121', '4c56ff4ce4aaf9573aa5dff913df9972', '121', 0, '', 1492858341, 0, 'offline', '/im/avatar/privateinvestocat.jpg', '80770675596'),
(6, 'simarui', '121', '87289608001', '司马睿', 0, '', 1493036784, 0, 'offline', '/im/avatar/twenty-percent-cooler-octocat.png', '73544123084'),
(7, 'simarui2', '121', '1850501660', 'simarui2', 0, '', 1493105196, 0, 'offline', '/im/avatar/twenty-percent-cooler-octocat.png', '78570265142');

-- --------------------------------------------------------

--
-- 表的结构 `user_join_group`
--

CREATE TABLE `user_join_group` (
  `id` int(11) NOT NULL,
  `uid` int(10) UNSIGNED NOT NULL,
  `group_id` int(10) UNSIGNED NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

--
-- 转存表中的数据 `user_join_group`
--

INSERT INTO `user_join_group` (`id`, `uid`, `group_id`) VALUES
(25, 4, 3),
(22, 7, 1),
(23, 7, 2),
(24, 7, 3),
(26, 7, 5);

--
-- Indexes for dumped tables
--

--
-- Indexes for table `contacts`
--
ALTER TABLE `contacts`
  ADD PRIMARY KEY (`id`),
  ADD KEY `master_uid` (`master_uid`),
  ADD KEY `master_uid_2` (`master_uid`,`group_id`);

--
-- Indexes for table `contact_group`
--
ALTER TABLE `contact_group`
  ADD PRIMARY KEY (`id`),
  ADD KEY `uid` (`uid`);

--
-- Indexes for table `global_group`
--
ALTER TABLE `global_group`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `channel_id` (`channel_id`);

--
-- Indexes for table `req_friend`
--
ALTER TABLE `req_friend`
  ADD PRIMARY KEY (`id`),
  ADD KEY `req_uid` (`req_uid`),
  ADD KEY `type` (`type`),
  ADD KEY `type_requid` (`type`,`req_uid`) USING BTREE,
  ADD KEY `time` (`time`);

--
-- Indexes for table `user`
--
ALTER TABLE `user`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `user` (`user`),
  ADD KEY `sid` (`sid`);

--
-- Indexes for table `user_join_group`
--
ALTER TABLE `user_join_group`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `uid_2` (`uid`,`group_id`),
  ADD KEY `uid` (`uid`);

--
-- 在导出的表使用AUTO_INCREMENT
--

--
-- 使用表AUTO_INCREMENT `contacts`
--
ALTER TABLE `contacts`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=8;
--
-- 使用表AUTO_INCREMENT `contact_group`
--
ALTER TABLE `contact_group`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=6;
--
-- 使用表AUTO_INCREMENT `global_group`
--
ALTER TABLE `global_group`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=6;
--
-- 使用表AUTO_INCREMENT `req_friend`
--
ALTER TABLE `req_friend`
  MODIFY `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=13;
--
-- 使用表AUTO_INCREMENT `user`
--
ALTER TABLE `user`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=8;
--
-- 使用表AUTO_INCREMENT `user_join_group`
--
ALTER TABLE `user_join_group`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=27;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
