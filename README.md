### 全网极速 VIP 视频助手

#####
```
apt -y update && apt -y install curl wget git unzip nginx mariadb-server vim

systemctl enable --now nginx mariadb

mysql初始化
mysql_secure_installation

登陆mysql创建数据库
mysql -u root -p
CREATE DATABASE video CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
GRANT ALL PRIVILEGES ON video.* TO video@localhost IDENTIFIED BY 'vip#video123!';
FLUSH PRIVILEGES;

创建数据表
CREATE TABLE `visitor_logs` (
  `id` INT AUTO_INCREMENT PRIMARY KEY,
  `ip_address` VARCHAR(45) NOT NULL,
  `visit_time` DATETIME DEFAULT CURRENT_TIMESTAMP,
  `user_agent` TEXT
);


quit

git clone https://github.com/hhttco/VipVideo.git
cd VipVideo
chown -R www-data:www-data /var/www/VipVideo
chmod -R 755 /var/www/VipVideo


go run main.go
```
