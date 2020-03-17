-- phpMyAdmin SQL Dump
-- version 5.0.0
-- https://www.phpmyadmin.net/
--
-- 主机： localhost
-- 生成日期： 2020-02-10 11:04:41
-- 服务器版本： 5.7.29-0ubuntu0.18.04.1
-- PHP 版本： 7.2.24-0ubuntu0.18.04.2

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
SET AUTOCOMMIT = 0;
START TRANSACTION;
SET time_zone = "+00:00";

/* create database massgrid default character set utf8 collate utf8_general_ci; */


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- 数据库： `massgrid`
--
CREATE DATABASE IF NOT EXISTS `massgrid` DEFAULT CHARACTER SET utf8 COLLATE utf8_bin;
USE `massgrid`;

-- --------------------------------------------------------

--
-- 表的结构 `articles`
--

DROP TABLE IF EXISTS `articles`;
CREATE TABLE `articles` (
  `id` int(10) UNSIGNED NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `article_id` int(11) DEFAULT NULL,
  `language` int(11) DEFAULT NULL,
  `pushed_at` bigint(20) DEFAULT NULL,
  `announcer` varchar(20) COLLATE utf8_bin DEFAULT NULL,
  `title` varchar(128) COLLATE utf8_bin DEFAULT NULL,
  `summary` varchar(128) COLLATE utf8_bin DEFAULT NULL,
  `content` text COLLATE utf8_bin
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

-- --------------------------------------------------------

--
-- 表的结构 `article_categories`
--

DROP TABLE IF EXISTS `article_categories`;
CREATE TABLE `article_categories` (
  `id` int(10) UNSIGNED NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `author_id` int(10) DEFAULT NULL,
  `type` int(2) DEFAULT NULL,
  `category` int(2) DEFAULT NULL,
  `name` varchar(20) COLLATE utf8_bin DEFAULT NULL,
  `language` int(2) DEFAULT NULL,
  `status` tinyint(1) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

-- --------------------------------------------------------

--
-- 表的结构 `article_infos`
--

DROP TABLE IF EXISTS `article_infos`;
CREATE TABLE `article_infos` (
  `id` int(10) UNSIGNED NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `article_id` int(11) DEFAULT NULL,
  `category` int(11) DEFAULT NULL,
  `author_id` int(11) DEFAULT NULL,
  `type` int(11) DEFAULT NULL,
  `priority` int(11) DEFAULT NULL,
  `status` tinyint(1) DEFAULT NULL,
  `read_count` int(11) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

-- --------------------------------------------------------

--
-- 表的结构 `assets`
--

DROP TABLE IF EXISTS `assets`;
CREATE TABLE `assets` (
  `id` int(10) UNSIGNED NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `flow_id` int(10) DEFAULT NULL,
  `coin_type` varchar(10) COLLATE utf8_bin DEFAULT NULL,
  `count` int(10) DEFAULT NULL,
  `increace_amount` varchar(30) COLLATE utf8_bin DEFAULT NULL,
  `reduce_amount` varchar(30) COLLATE utf8_bin DEFAULT NULL,
  `total_amount` varchar(30) COLLATE utf8_bin DEFAULT NULL,
  `start_at` bigint(20) DEFAULT NULL,
  `end_at` bigint(20) DEFAULT NULL,
  `total_time` bigint(20) DEFAULT NULL,
  `author_id` int(11) DEFAULT NULL,
  `create_at` bigint(20) DEFAULT NULL,
  `status` tinyint(1) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

-- --------------------------------------------------------

--
-- 表的结构 `btc_txes`
--

DROP TABLE IF EXISTS `btc_txes`;
CREATE TABLE `btc_txes` (
  `id` int(10) UNSIGNED NOT NULL,
  `tx_index` int(10) DEFAULT NULL,
  `tx_id` varchar(80) COLLATE utf8_bin DEFAULT NULL,
  `to_address` varchar(80) COLLATE utf8_bin DEFAULT NULL,
  `tx_type` varchar(10) COLLATE utf8_bin DEFAULT NULL,
  `to_amount` varchar(30) COLLATE utf8_bin DEFAULT NULL,
  `fee` varchar(30) COLLATE utf8_bin DEFAULT NULL,
  `time` bigint(20) DEFAULT NULL,
  `server_ip` varchar(50) COLLATE utf8_bin DEFAULT NULL,
  `detail` varchar(4500) COLLATE utf8_bin DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

-- --------------------------------------------------------

--
-- 表的结构 `coin_markets`
--

DROP TABLE IF EXISTS `coin_markets`;
CREATE TABLE `coin_markets` (
  `id` int(10) UNSIGNED NOT NULL,
  `name` varchar(12) COLLATE utf8_bin DEFAULT NULL,
  `url` varchar(32) COLLATE utf8_bin DEFAULT NULL,
  `price` varchar(32) COLLATE utf8_bin DEFAULT NULL,
  `time` bigint(20) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

-- --------------------------------------------------------

--
-- 表的结构 `coin_prices`
--

DROP TABLE IF EXISTS `coin_prices`;
CREATE TABLE `coin_prices` (
  `id` int(10) UNSIGNED NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `name` varchar(12) COLLATE utf8_bin DEFAULT NULL,
  `price` varchar(32) COLLATE utf8_bin DEFAULT NULL,
  `discount` varchar(12) COLLATE utf8_bin DEFAULT NULL,
  `auto_update` tinyint(1) DEFAULT NULL,
  `time` bigint(20) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

-- --------------------------------------------------------

--
-- 表的结构 `eth_txes`
--

DROP TABLE IF EXISTS `eth_txes`;
CREATE TABLE `eth_txes` (
  `id` int(10) UNSIGNED NOT NULL,
  `block_height` bigint(20) DEFAULT NULL,
  `tx_hash` varchar(70) COLLATE utf8_bin DEFAULT NULL,
  `tx_index` int(10) DEFAULT NULL,
  `coin_from` varchar(64) COLLATE utf8_bin DEFAULT NULL,
  `coin_to` varchar(64) COLLATE utf8_bin DEFAULT NULL,
  `tx_type` varchar(64) COLLATE utf8_bin DEFAULT NULL,
  `nonce` int(12) DEFAULT NULL,
  `value` varchar(30) COLLATE utf8_bin DEFAULT NULL,
  `gas` int(10) DEFAULT NULL,
  `gas_price` varchar(30) COLLATE utf8_bin DEFAULT NULL,
  `server_ip` varchar(255) COLLATE utf8_bin DEFAULT NULL,
  `time` bigint(20) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

-- --------------------------------------------------------

--
-- 表的结构 `farm_servers`
--

DROP TABLE IF EXISTS `farm_servers`;
CREATE TABLE `farm_servers` (
  `id` int(10) UNSIGNED NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `farm_id` varchar(64) COLLATE utf8_bin DEFAULT NULL,
  `miner_type` varchar(32) COLLATE utf8_bin DEFAULT NULL,
  `available_count` bigint(20) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

-- --------------------------------------------------------

--
-- 表的结构 `mgd_txes`
--

DROP TABLE IF EXISTS `mgd_txes`;
CREATE TABLE `mgd_txes` (
  `id` int(10) UNSIGNED NOT NULL,
  `tx_index` int(10) DEFAULT NULL,
  `tx_id` varchar(80) COLLATE utf8_bin DEFAULT NULL,
  `to_address` varchar(80) COLLATE utf8_bin DEFAULT NULL,
  `tx_type` varchar(10) COLLATE utf8_bin DEFAULT NULL,
  `to_amount` varchar(30) COLLATE utf8_bin DEFAULT NULL,
  `fee` varchar(30) COLLATE utf8_bin DEFAULT NULL,
  `time` bigint(20) DEFAULT NULL,
  `server_ip` varchar(50) COLLATE utf8_bin DEFAULT NULL,
  `detail` varchar(4500) COLLATE utf8_bin DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

-- --------------------------------------------------------

--
-- 表的结构 `miner_accounts`
--

DROP TABLE IF EXISTS `miner_accounts`;
CREATE TABLE `miner_accounts` (
  `id` int(10) UNSIGNED NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `user_id` int(10) DEFAULT NULL,
  `coin_type` varchar(10) COLLATE utf8_bin DEFAULT NULL,
  `coin_addr` varchar(80) COLLATE utf8_bin DEFAULT NULL,
  `addr_type` int(10) DEFAULT NULL,
  `coin_amount` varchar(30) COLLATE utf8_bin DEFAULT NULL,
  `total_amount` varchar(30) COLLATE utf8_bin DEFAULT NULL,
  `status` int(10) DEFAULT NULL,
  `miner_pool` varchar(80) COLLATE utf8_bin DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

-- --------------------------------------------------------

--
-- 表的结构 `miner_orders`
--

DROP TABLE IF EXISTS `miner_orders`;
CREATE TABLE `miner_orders` (
  `id` int(10) UNSIGNED NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `goods_id` varchar(30) COLLATE utf8_bin DEFAULT NULL,
  `miner_order_id` varchar(64) COLLATE utf8_bin DEFAULT NULL,
  `farm_id` varchar(64) COLLATE utf8_bin DEFAULT NULL,
  `miner_id` varchar(64) COLLATE utf8_bin DEFAULT NULL,
  `create_at` bigint(20) DEFAULT NULL,
  `update_at` bigint(20) DEFAULT NULL,
  `status` int(5) DEFAULT NULL,
  `goods_type` varchar(32) COLLATE utf8_bin DEFAULT NULL,
  `goods_price` bigint(20) DEFAULT NULL,
  `rent_time` bigint(20) DEFAULT NULL,
  `miner_pool` varchar(128) COLLATE utf8_bin DEFAULT NULL,
  `miner_username` varchar(64) COLLATE utf8_bin DEFAULT NULL,
  `miner_worker` varchar(64) COLLATE utf8_bin DEFAULT NULL,
  `out_trade_no` varchar(80) COLLATE utf8_bin DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

-- --------------------------------------------------------

--
-- 表的结构 `orders`
--

DROP TABLE IF EXISTS `orders`;
CREATE TABLE `orders` (
  `id` int(10) UNSIGNED NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `user_id` int(10) DEFAULT NULL,
  `pay_type` int(11) DEFAULT NULL,
  `pay_amount` varchar(11) COLLATE utf8_bin DEFAULT NULL,
  `ex_price` varchar(11) COLLATE utf8_bin DEFAULT NULL,
  `order_type` int(11) DEFAULT NULL,
  `trade_type` int(11) DEFAULT NULL,
  `trade_no` varchar(64) COLLATE utf8_bin DEFAULT NULL,
  `out_trade_no` varchar(64) COLLATE utf8_bin DEFAULT NULL,
  `total_amount` varchar(11) COLLATE utf8_bin DEFAULT NULL,
  `total_time` bigint(20) DEFAULT NULL,
  `create_at` bigint(20) DEFAULT NULL,
  `order_subject` varchar(20) COLLATE utf8_bin DEFAULT NULL,
  `order_status` int(10) DEFAULT NULL,
  `order_detail` text COLLATE utf8_bin,
  `rental_type` int(11) DEFAULT NULL,
  `goods_id` varchar(32) COLLATE utf8_bin DEFAULT NULL,
  `goods_name` varchar(32) COLLATE utf8_bin DEFAULT NULL,
  `goods_price` varchar(30) COLLATE utf8_bin DEFAULT NULL,
  `goods_quantity` int(10) DEFAULT NULL,
  `goods_unit` varchar(10) COLLATE utf8_bin DEFAULT NULL,
  `miner_goods_type` varchar(32) COLLATE utf8_bin DEFAULT NULL,
  `farm_id` varchar(32) COLLATE utf8_bin DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

-- --------------------------------------------------------

--
-- 表的结构 `pay_txes`
--

DROP TABLE IF EXISTS `pay_txes`;
CREATE TABLE `pay_txes` (
  `id` int(10) UNSIGNED NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `pay_channel` int(11) DEFAULT NULL,
  `pay_type` int(11) DEFAULT NULL,
  `trade_type` int(11) DEFAULT NULL,
  `order_id` varchar(20) COLLATE utf8_bin DEFAULT NULL,
  `user_id` int(10) DEFAULT NULL,
  `trade_no` varchar(64) COLLATE utf8_bin DEFAULT NULL,
  `out_trade_no` varchar(64) COLLATE utf8_bin DEFAULT NULL,
  `total_amount` varchar(11) COLLATE utf8_bin DEFAULT NULL,
  `remain_amount` varchar(11) COLLATE utf8_bin DEFAULT NULL,
  `recv_time` bigint(20) DEFAULT NULL,
  `trade_status` int(10) DEFAULT NULL,
  `buyer_user_id` varchar(16) COLLATE utf8_bin DEFAULT NULL,
  `buyer_logon_id` varchar(64) COLLATE utf8_bin DEFAULT NULL,
  `send_pay_date` varchar(32) COLLATE utf8_bin DEFAULT NULL,
  `order_descript` varchar(20) COLLATE utf8_bin DEFAULT NULL,
  `trade_detail` text COLLATE utf8_bin
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

-- --------------------------------------------------------

--
-- 表的结构 `products`
--

DROP TABLE IF EXISTS `products`;
CREATE TABLE `products` (
  `id` int(10) UNSIGNED NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `author_id` int(11) DEFAULT NULL,
  `goods_type` int(11) DEFAULT NULL,
  `goods_id` varchar(32) COLLATE utf8_bin DEFAULT NULL,
  `goods_name` varchar(32) COLLATE utf8_bin DEFAULT NULL,
  `org_price` varchar(30) COLLATE utf8_bin DEFAULT NULL,
  `cur_price` varchar(30) COLLATE utf8_bin DEFAULT NULL,
  `quantity` int(11) DEFAULT NULL,
  `total_quantity` int(11) DEFAULT NULL,
  `unit` varchar(10) COLLATE utf8_bin DEFAULT NULL,
  `description` varchar(30) COLLATE utf8_bin DEFAULT NULL,
  `image_uri` varchar(128) COLLATE utf8_bin DEFAULT NULL,
  `pushed_at` bigint(20) DEFAULT NULL,
  `status` tinyint(1) DEFAULT NULL,
  `farm_id` varchar(64) COLLATE utf8_bin DEFAULT NULL,
  `miner_goods_type` varchar(64) COLLATE utf8_bin DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

-- --------------------------------------------------------

--
-- 表的结构 `product_details`
--

DROP TABLE IF EXISTS `product_details`;
CREATE TABLE `product_details` (
  `id` int(10) UNSIGNED NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `products_id` int(10) DEFAULT NULL,
  `label` int(11) DEFAULT NULL,
  `coin_type` varchar(10) COLLATE utf8_bin DEFAULT NULL,
  `roi_st` int(10) DEFAULT NULL,
  `roi_end` int(10) DEFAULT NULL,
  `mining` int(10) DEFAULT NULL,
  `mining_unit` varchar(10) COLLATE utf8_bin DEFAULT NULL,
  `desirable_output` varchar(32) COLLATE utf8_bin DEFAULT NULL,
  `power` varchar(20) COLLATE utf8_bin DEFAULT NULL,
  `power_price` varchar(20) COLLATE utf8_bin DEFAULT NULL,
  `occupy_price` varchar(20) COLLATE utf8_bin DEFAULT NULL,
  `manger_price` varchar(20) COLLATE utf8_bin DEFAULT NULL,
  `lease_time` varchar(20) COLLATE utf8_bin DEFAULT NULL,
  `start_time` varchar(20) COLLATE utf8_bin DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

-- --------------------------------------------------------

--
-- 表的结构 `settings`
--

DROP TABLE IF EXISTS `settings`;
CREATE TABLE `settings` (
  `id` int(10) UNSIGNED NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `author_id` int(10) DEFAULT NULL,
  `category` varchar(32) COLLATE utf8_bin DEFAULT NULL,
  `name` varchar(32) COLLATE utf8_bin DEFAULT NULL,
  `value` text COLLATE utf8_bin
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

-- --------------------------------------------------------

--
-- 表的结构 `system_accounts`
--

DROP TABLE IF EXISTS `system_accounts`;
CREATE TABLE `system_accounts` (
  `id` int(10) UNSIGNED NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `coin_type` varchar(10) COLLATE utf8_bin DEFAULT NULL,
  `coin_amount` varchar(30) COLLATE utf8_bin DEFAULT NULL,
  `update_at` bigint(20) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

-- --------------------------------------------------------

--
-- 表的结构 `system_basics`
--

DROP TABLE IF EXISTS `system_basics`;
CREATE TABLE `system_basics` (
  `id` int(10) UNSIGNED NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `uid` varchar(20) COLLATE utf8_bin DEFAULT NULL,
  `nick_name` varchar(64) COLLATE utf8_bin DEFAULT NULL,
  `net_remote_ip` varchar(20) COLLATE utf8_bin DEFAULT NULL,
  `cpu_count` int(11) DEFAULT NULL,
  `cpu_name` varchar(64) COLLATE utf8_bin DEFAULT NULL,
  `os_arch` varchar(20) COLLATE utf8_bin DEFAULT NULL,
  `os_byte_order` varchar(20) COLLATE utf8_bin DEFAULT NULL,
  `os_system` varchar(20) COLLATE utf8_bin DEFAULT NULL,
  `net_mac` varchar(20) COLLATE utf8_bin DEFAULT NULL,
  `disk_total` varchar(20) COLLATE utf8_bin DEFAULT NULL,
  `mem_total` varchar(20) COLLATE utf8_bin DEFAULT NULL,
  `update_at` bigint(20) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

-- --------------------------------------------------------

--
-- 表的结构 `system_simples`
--

DROP TABLE IF EXISTS `system_simples`;
CREATE TABLE `system_simples` (
  `id` int(10) UNSIGNED NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `uid` varchar(20) COLLATE utf8_bin DEFAULT NULL,
  `net_remote_ip` varchar(20) COLLATE utf8_bin DEFAULT NULL,
  `net_local_ip` varchar(20) COLLATE utf8_bin DEFAULT NULL,
  `mem_used` varchar(20) COLLATE utf8_bin DEFAULT NULL,
  `mem_usage` varchar(12) COLLATE utf8_bin DEFAULT NULL,
  `sys_uptime` varchar(64) COLLATE utf8_bin DEFAULT NULL,
  `disk_used` varchar(20) COLLATE utf8_bin DEFAULT NULL,
  `disk_usage` varchar(12) COLLATE utf8_bin DEFAULT NULL,
  `net_byte_sent` bigint(20) DEFAULT NULL,
  `net_byte_recv` bigint(20) DEFAULT NULL,
  `cpu_average` varchar(128) COLLATE utf8_bin DEFAULT NULL,
  `login_count` int(11) DEFAULT NULL,
  `login_user` varchar(1024) COLLATE utf8_bin DEFAULT NULL,
  `update_at` bigint(20) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

-- --------------------------------------------------------

--
-- 表的结构 `users`
--

DROP TABLE IF EXISTS `users`;
CREATE TABLE `users` (
  `id` int(10) UNSIGNED NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `user_id` int(11) DEFAULT NULL,
  `username` varchar(30) COLLATE utf8_bin DEFAULT NULL,
  `email` varchar(20) COLLATE utf8_bin DEFAULT NULL,
  `password` varchar(255) COLLATE utf8_bin DEFAULT NULL,
  `phone` varchar(20) COLLATE utf8_bin DEFAULT NULL,
  `role_id` int(11) DEFAULT NULL,
  `register_at` bigint(20) DEFAULT NULL,
  `register_ip` varchar(50) COLLATE utf8_bin DEFAULT NULL,
  `user_status` bigint(20) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

--
-- 转存表中的数据 `users`
--

INSERT INTO `users` (`id`, `created_at`, `updated_at`, `deleted_at`, `user_id`, `username`, `email`, `password`, `phone`, `role_id`, `register_at`, `register_ip`, `user_status`) VALUES
(1, '2019-12-04 01:59:39', '2019-12-04 01:59:39', NULL, 1, 'zhxx123', 'zhxx_123@qq.com', '$2a$10$NLoeCbKEXA9ax76DqWTsF.EwGxEzduo61Nk.XY0DbfYr6dhp0R4SK', '17168970023', 5, 1575424779, '', 0);

-- --------------------------------------------------------

--
-- 表的结构 `user_accounts`
--

DROP TABLE IF EXISTS `user_accounts`;
CREATE TABLE `user_accounts` (
  `id` int(10) UNSIGNED NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `user_id` int(11) DEFAULT NULL,
  `coin_type` varchar(10) COLLATE utf8_bin DEFAULT NULL,
  `coin_addr` varchar(80) COLLATE utf8_bin NOT NULL,
  `coin_amount` varchar(30) COLLATE utf8_bin DEFAULT NULL,
  `virtual_amount` varchar(30) COLLATE utf8_bin DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

-- --------------------------------------------------------

--
-- 表的结构 `user_assetflows`
--

DROP TABLE IF EXISTS `user_assetflows`;
CREATE TABLE `user_assetflows` (
  `id` int(10) UNSIGNED NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `user_id` int(10) DEFAULT NULL,
  `out_trade_no` varchar(64) COLLATE utf8_bin DEFAULT NULL,
  `trade_type` int(11) DEFAULT NULL,
  `create_at` bigint(20) DEFAULT NULL,
  `coin_type` varchar(10) COLLATE utf8_bin DEFAULT NULL,
  `amount` varchar(30) COLLATE utf8_bin DEFAULT NULL,
  `total_amount` varchar(30) COLLATE utf8_bin DEFAULT NULL,
  `description` varchar(64) COLLATE utf8_bin DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

-- --------------------------------------------------------

--
-- 表的结构 `user_login_sets`
--

DROP TABLE IF EXISTS `user_login_sets`;
CREATE TABLE `user_login_sets` (
  `id` int(10) UNSIGNED NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `user_id` int(11) DEFAULT NULL,
  `login_count` int(11) DEFAULT NULL,
  `login_at` bigint(20) DEFAULT NULL,
  `ip_list` varchar(255) COLLATE utf8_bin DEFAULT NULL,
  `verify` tinyint(1) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

-- --------------------------------------------------------

--
-- 表的结构 `user_messages`
--

DROP TABLE IF EXISTS `user_messages`;
CREATE TABLE `user_messages` (
  `id` int(10) UNSIGNED NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `user_id` int(11) DEFAULT NULL,
  `article_id` int(11) DEFAULT NULL,
  `category` int(11) DEFAULT NULL,
  `title` varchar(128) COLLATE utf8_bin DEFAULT NULL,
  `summary` varchar(128) COLLATE utf8_bin DEFAULT NULL,
  `content` text COLLATE utf8_bin,
  `pushed_at` bigint(20) DEFAULT NULL,
  `announcer` varchar(20) COLLATE utf8_bin DEFAULT NULL,
  `author_id` int(11) DEFAULT NULL,
  `readed` tinyint(1) DEFAULT NULL,
  `read_at` bigint(20) DEFAULT NULL,
  `status` tinyint(1) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

-- --------------------------------------------------------

--
-- 表的结构 `user_oauths`
--

DROP TABLE IF EXISTS `user_oauths`;
CREATE TABLE `user_oauths` (
  `id` int(10) UNSIGNED NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `user_id` int(11) DEFAULT NULL,
  `role_id` int(11) DEFAULT NULL,
  `login_at` bigint(20) DEFAULT NULL,
  `token` varchar(255) COLLATE utf8_bin DEFAULT NULL,
  `secret` varchar(255) COLLATE utf8_bin DEFAULT NULL,
  `login_ip` varchar(50) COLLATE utf8_bin DEFAULT NULL,
  `login_city` varchar(50) COLLATE utf8_bin DEFAULT NULL,
  `express_in` bigint(20) DEFAULT NULL,
  `revoked` tinyint(1) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

-- --------------------------------------------------------

--
-- 表的结构 `virtual_accounts`
--

DROP TABLE IF EXISTS `virtual_accounts`;
CREATE TABLE `virtual_accounts` (
  `id` int(10) UNSIGNED NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `user_id` int(10) DEFAULT NULL,
  `out_trade_no` varchar(64) COLLATE utf8_bin DEFAULT NULL,
  `create_at` bigint(20) DEFAULT NULL,
  `coin_type` varchar(10) COLLATE utf8_bin DEFAULT NULL,
  `coin_amount` varchar(30) COLLATE utf8_bin DEFAULT NULL,
  `operator_id` int(10) DEFAULT NULL,
  `recharge_amount` varchar(30) COLLATE utf8_bin DEFAULT NULL,
  `description` varchar(64) COLLATE utf8_bin DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

-- --------------------------------------------------------

--
-- 表的结构 `virtual_recharges`
--

DROP TABLE IF EXISTS `virtual_recharges`;
CREATE TABLE `virtual_recharges` (
  `id` int(10) UNSIGNED NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `user_id` int(10) DEFAULT NULL,
  `out_trade_no` varchar(64) COLLATE utf8_bin DEFAULT NULL,
  `create_at` bigint(20) DEFAULT NULL,
  `coin_type` varchar(10) COLLATE utf8_bin DEFAULT NULL,
  `coin_amount` varchar(30) COLLATE utf8_bin DEFAULT NULL,
  `operator_id` int(10) DEFAULT NULL,
  `recharge_amount` varchar(30) COLLATE utf8_bin DEFAULT NULL,
  `description` varchar(64) COLLATE utf8_bin DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

-- --------------------------------------------------------

--
-- 表的结构 `wallet_addresses`
--

DROP TABLE IF EXISTS `wallet_addresses`;
CREATE TABLE `wallet_addresses` (
  `id` int(10) UNSIGNED NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `indexs` int(10) DEFAULT NULL,
  `coin_type` varchar(255) COLLATE utf8_bin DEFAULT NULL,
  `account` varchar(255) COLLATE utf8_bin DEFAULT NULL,
  `address` varchar(80) COLLATE utf8_bin DEFAULT NULL,
  `user_id` int(10) DEFAULT NULL,
  `allocated` tinyint(1) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

-- --------------------------------------------------------

--
-- 表的结构 `wallet_basics`
--

DROP TABLE IF EXISTS `wallet_basics`;
CREATE TABLE `wallet_basics` (
  `id` int(10) UNSIGNED NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `name` varchar(20) COLLATE utf8_bin DEFAULT NULL,
  `version` varchar(10) COLLATE utf8_bin DEFAULT NULL,
  `net_model` varchar(255) COLLATE utf8_bin DEFAULT NULL,
  `update_at` bigint(20) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

-- --------------------------------------------------------

--
-- 表的结构 `wallet_records`
--

DROP TABLE IF EXISTS `wallet_records`;
CREATE TABLE `wallet_records` (
  `id` int(10) UNSIGNED NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `tx_id` varchar(255) COLLATE utf8_bin DEFAULT NULL,
  `coin_type` varchar(255) COLLATE utf8_bin DEFAULT NULL,
  `added` tinyint(1) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

-- --------------------------------------------------------

--
-- 表的结构 `wallet_simples`
--

DROP TABLE IF EXISTS `wallet_simples`;
CREATE TABLE `wallet_simples` (
  `id` int(10) UNSIGNED NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `name` varchar(20) COLLATE utf8_bin DEFAULT NULL,
  `balance` varchar(32) COLLATE utf8_bin DEFAULT NULL,
  `block_height` bigint(20) DEFAULT NULL,
  `difficulty` varchar(32) COLLATE utf8_bin DEFAULT NULL,
  `network_hash` varchar(32) COLLATE utf8_bin DEFAULT NULL,
  `update_at` bigint(20) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

-- --------------------------------------------------------

--
-- 表的结构 `wallet_syncs`
--

DROP TABLE IF EXISTS `wallet_syncs`;
CREATE TABLE `wallet_syncs` (
  `id` int(10) UNSIGNED NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `coin_type` varchar(255) COLLATE utf8_bin DEFAULT NULL,
  `start_block` bigint(20) DEFAULT NULL,
  `last_block` bigint(20) DEFAULT NULL,
  `update_at` bigint(20) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

-- --------------------------------------------------------

--
-- 表的结构 `wallet_txes`
--

DROP TABLE IF EXISTS `wallet_txes`;
CREATE TABLE `wallet_txes` (
  `id` int(10) UNSIGNED NOT NULL,
  `tx_index` int(10) DEFAULT NULL,
  `tx_id` varchar(80) COLLATE utf8_bin DEFAULT NULL,
  `to_address` varchar(80) COLLATE utf8_bin DEFAULT NULL,
  `tx_type` varchar(10) COLLATE utf8_bin DEFAULT NULL,
  `to_amount` varchar(30) COLLATE utf8_bin DEFAULT NULL,
  `fee` varchar(30) COLLATE utf8_bin DEFAULT NULL,
  `time` bigint(20) DEFAULT NULL,
  `server_ip` varchar(50) COLLATE utf8_bin DEFAULT NULL,
  `detail` varchar(4500) COLLATE utf8_bin DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

-- --------------------------------------------------------

--
-- 表的结构 `work_orders`
--

DROP TABLE IF EXISTS `work_orders`;
CREATE TABLE `work_orders` (
  `id` int(10) UNSIGNED NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `user_id` int(11) DEFAULT NULL,
  `work_id` int(11) DEFAULT NULL,
  `issue_type` varchar(10) COLLATE utf8_bin DEFAULT NULL,
  `description` varchar(255) COLLATE utf8_bin DEFAULT NULL,
  `email` varchar(50) COLLATE utf8_bin DEFAULT NULL,
  `img_uri` varchar(128) COLLATE utf8_bin DEFAULT NULL,
  `create_at` bigint(20) DEFAULT NULL,
  `status` varchar(10) COLLATE utf8_bin DEFAULT NULL,
  `operator_id` int(11) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

--
-- 转储表的索引
--

--
-- 表的索引 `articles`
--
ALTER TABLE `articles`
  ADD PRIMARY KEY (`id`),
  ADD KEY `idx_articles_deleted_at` (`deleted_at`);

--
-- 表的索引 `article_categories`
--
ALTER TABLE `article_categories`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `category` (`category`),
  ADD KEY `idx_article_categories_deleted_at` (`deleted_at`);

--
-- 表的索引 `article_infos`
--
ALTER TABLE `article_infos`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `article_id` (`article_id`),
  ADD KEY `idx_article_infos_deleted_at` (`deleted_at`);

--
-- 表的索引 `assets`
--
ALTER TABLE `assets`
  ADD PRIMARY KEY (`id`),
  ADD KEY `idx_assets_deleted_at` (`deleted_at`);

--
-- 表的索引 `btc_txes`
--
ALTER TABLE `btc_txes`
  ADD PRIMARY KEY (`id`);

--
-- 表的索引 `coin_markets`
--
ALTER TABLE `coin_markets`
  ADD PRIMARY KEY (`id`);

--
-- 表的索引 `coin_prices`
--
ALTER TABLE `coin_prices`
  ADD PRIMARY KEY (`id`),
  ADD KEY `idx_coin_prices_deleted_at` (`deleted_at`);

--
-- 表的索引 `eth_txes`
--
ALTER TABLE `eth_txes`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `tx_hash` (`tx_hash`);

--
-- 表的索引 `farm_servers`
--
ALTER TABLE `farm_servers`
  ADD PRIMARY KEY (`id`),
  ADD KEY `idx_farm_servers_deleted_at` (`deleted_at`);

--
-- 表的索引 `mgd_txes`
--
ALTER TABLE `mgd_txes`
  ADD PRIMARY KEY (`id`);

--
-- 表的索引 `miner_accounts`
--
ALTER TABLE `miner_accounts`
  ADD PRIMARY KEY (`id`),
  ADD KEY `idx_miner_accounts_deleted_at` (`deleted_at`);

--
-- 表的索引 `miner_orders`
--
ALTER TABLE `miner_orders`
  ADD PRIMARY KEY (`id`),
  ADD KEY `idx_miner_orders_deleted_at` (`deleted_at`);

--
-- 表的索引 `orders`
--
ALTER TABLE `orders`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `trade_no` (`trade_no`),
  ADD KEY `idx_orders_deleted_at` (`deleted_at`);

--
-- 表的索引 `pay_txes`
--
ALTER TABLE `pay_txes`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `trade_no` (`trade_no`),
  ADD KEY `idx_pay_txes_deleted_at` (`deleted_at`);

--
-- 表的索引 `products`
--
ALTER TABLE `products`
  ADD PRIMARY KEY (`id`),
  ADD KEY `idx_products_deleted_at` (`deleted_at`);

--
-- 表的索引 `product_details`
--
ALTER TABLE `product_details`
  ADD PRIMARY KEY (`id`),
  ADD KEY `idx_product_details_deleted_at` (`deleted_at`);

--
-- 表的索引 `settings`
--
ALTER TABLE `settings`
  ADD PRIMARY KEY (`id`),
  ADD KEY `idx_settings_deleted_at` (`deleted_at`);

--
-- 表的索引 `system_accounts`
--
ALTER TABLE `system_accounts`
  ADD PRIMARY KEY (`id`),
  ADD KEY `idx_system_accounts_deleted_at` (`deleted_at`);

--
-- 表的索引 `system_basics`
--
ALTER TABLE `system_basics`
  ADD PRIMARY KEY (`id`),
  ADD KEY `idx_system_basics_deleted_at` (`deleted_at`);

--
-- 表的索引 `system_simples`
--
ALTER TABLE `system_simples`
  ADD PRIMARY KEY (`id`),
  ADD KEY `idx_system_simples_deleted_at` (`deleted_at`);

--
-- 表的索引 `users`
--
ALTER TABLE `users`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `user_id` (`user_id`),
  ADD UNIQUE KEY `email` (`email`),
  ADD KEY `idx_users_deleted_at` (`deleted_at`);

--
-- 表的索引 `user_accounts`
--
ALTER TABLE `user_accounts`
  ADD PRIMARY KEY (`id`),
  ADD KEY `idx_user_accounts_deleted_at` (`deleted_at`);

--
-- 表的索引 `user_assetflows`
--
ALTER TABLE `user_assetflows`
  ADD PRIMARY KEY (`id`),
  ADD KEY `idx_user_assetflows_deleted_at` (`deleted_at`);

--
-- 表的索引 `user_login_sets`
--
ALTER TABLE `user_login_sets`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `login_count` (`login_count`),
  ADD KEY `idx_user_login_sets_deleted_at` (`deleted_at`);

--
-- 表的索引 `user_messages`
--
ALTER TABLE `user_messages`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `article_id` (`article_id`),
  ADD KEY `idx_user_messages_deleted_at` (`deleted_at`);

--
-- 表的索引 `user_oauths`
--
ALTER TABLE `user_oauths`
  ADD PRIMARY KEY (`id`),
  ADD KEY `idx_user_oauths_deleted_at` (`deleted_at`);

--
-- 表的索引 `virtual_accounts`
--
ALTER TABLE `virtual_accounts`
  ADD PRIMARY KEY (`id`),
  ADD KEY `idx_virtual_accounts_deleted_at` (`deleted_at`);

--
-- 表的索引 `virtual_recharges`
--
ALTER TABLE `virtual_recharges`
  ADD PRIMARY KEY (`id`),
  ADD KEY `idx_virtual_recharges_deleted_at` (`deleted_at`);

--
-- 表的索引 `wallet_addresses`
--
ALTER TABLE `wallet_addresses`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `address` (`address`),
  ADD KEY `idx_wallet_addresses_deleted_at` (`deleted_at`);

--
-- 表的索引 `wallet_basics`
--
ALTER TABLE `wallet_basics`
  ADD PRIMARY KEY (`id`),
  ADD KEY `idx_wallet_basics_deleted_at` (`deleted_at`);

--
-- 表的索引 `wallet_records`
--
ALTER TABLE `wallet_records`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `tx_id` (`tx_id`),
  ADD KEY `idx_wallet_records_deleted_at` (`deleted_at`);

--
-- 表的索引 `wallet_simples`
--
ALTER TABLE `wallet_simples`
  ADD PRIMARY KEY (`id`),
  ADD KEY `idx_wallet_simples_deleted_at` (`deleted_at`);

--
-- 表的索引 `wallet_syncs`
--
ALTER TABLE `wallet_syncs`
  ADD PRIMARY KEY (`id`),
  ADD KEY `idx_wallet_syncs_deleted_at` (`deleted_at`);

--
-- 表的索引 `wallet_txes`
--
ALTER TABLE `wallet_txes`
  ADD PRIMARY KEY (`id`);

--
-- 表的索引 `work_orders`
--
ALTER TABLE `work_orders`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `work_id` (`work_id`),
  ADD UNIQUE KEY `operator_id` (`operator_id`),
  ADD KEY `idx_work_orders_deleted_at` (`deleted_at`);

--
-- 在导出的表使用AUTO_INCREMENT
--

--
-- 使用表AUTO_INCREMENT `articles`
--
ALTER TABLE `articles`
  MODIFY `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- 使用表AUTO_INCREMENT `article_categories`
--
ALTER TABLE `article_categories`
  MODIFY `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- 使用表AUTO_INCREMENT `article_infos`
--
ALTER TABLE `article_infos`
  MODIFY `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- 使用表AUTO_INCREMENT `assets`
--
ALTER TABLE `assets`
  MODIFY `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- 使用表AUTO_INCREMENT `btc_txes`
--
ALTER TABLE `btc_txes`
  MODIFY `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- 使用表AUTO_INCREMENT `coin_markets`
--
ALTER TABLE `coin_markets`
  MODIFY `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- 使用表AUTO_INCREMENT `coin_prices`
--
ALTER TABLE `coin_prices`
  MODIFY `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- 使用表AUTO_INCREMENT `eth_txes`
--
ALTER TABLE `eth_txes`
  MODIFY `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- 使用表AUTO_INCREMENT `farm_servers`
--
ALTER TABLE `farm_servers`
  MODIFY `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- 使用表AUTO_INCREMENT `mgd_txes`
--
ALTER TABLE `mgd_txes`
  MODIFY `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- 使用表AUTO_INCREMENT `miner_accounts`
--
ALTER TABLE `miner_accounts`
  MODIFY `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- 使用表AUTO_INCREMENT `miner_orders`
--
ALTER TABLE `miner_orders`
  MODIFY `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- 使用表AUTO_INCREMENT `orders`
--
ALTER TABLE `orders`
  MODIFY `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- 使用表AUTO_INCREMENT `pay_txes`
--
ALTER TABLE `pay_txes`
  MODIFY `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- 使用表AUTO_INCREMENT `products`
--
ALTER TABLE `products`
  MODIFY `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- 使用表AUTO_INCREMENT `product_details`
--
ALTER TABLE `product_details`
  MODIFY `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- 使用表AUTO_INCREMENT `settings`
--
ALTER TABLE `settings`
  MODIFY `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- 使用表AUTO_INCREMENT `system_accounts`
--
ALTER TABLE `system_accounts`
  MODIFY `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- 使用表AUTO_INCREMENT `system_basics`
--
ALTER TABLE `system_basics`
  MODIFY `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- 使用表AUTO_INCREMENT `system_simples`
--
ALTER TABLE `system_simples`
  MODIFY `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- 使用表AUTO_INCREMENT `users`
--
ALTER TABLE `users`
  MODIFY `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=16;

--
-- 使用表AUTO_INCREMENT `user_accounts`
--
ALTER TABLE `user_accounts`
  MODIFY `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- 使用表AUTO_INCREMENT `user_assetflows`
--
ALTER TABLE `user_assetflows`
  MODIFY `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- 使用表AUTO_INCREMENT `user_login_sets`
--
ALTER TABLE `user_login_sets`
  MODIFY `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- 使用表AUTO_INCREMENT `user_messages`
--
ALTER TABLE `user_messages`
  MODIFY `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- 使用表AUTO_INCREMENT `user_oauths`
--
ALTER TABLE `user_oauths`
  MODIFY `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- 使用表AUTO_INCREMENT `virtual_accounts`
--
ALTER TABLE `virtual_accounts`
  MODIFY `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- 使用表AUTO_INCREMENT `virtual_recharges`
--
ALTER TABLE `virtual_recharges`
  MODIFY `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- 使用表AUTO_INCREMENT `wallet_addresses`
--
ALTER TABLE `wallet_addresses`
  MODIFY `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- 使用表AUTO_INCREMENT `wallet_basics`
--
ALTER TABLE `wallet_basics`
  MODIFY `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- 使用表AUTO_INCREMENT `wallet_records`
--
ALTER TABLE `wallet_records`
  MODIFY `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- 使用表AUTO_INCREMENT `wallet_simples`
--
ALTER TABLE `wallet_simples`
  MODIFY `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- 使用表AUTO_INCREMENT `wallet_syncs`
--
ALTER TABLE `wallet_syncs`
  MODIFY `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- 使用表AUTO_INCREMENT `wallet_txes`
--
ALTER TABLE `wallet_txes`
  MODIFY `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- 使用表AUTO_INCREMENT `work_orders`
--
ALTER TABLE `work_orders`
  MODIFY `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT;
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
