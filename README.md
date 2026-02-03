### 全网极速 VIP 视频助手

### debin 10 + 部署
```
apt -y update && apt -y install curl wget git unzip nginx mariadb-server vim
```

设置开机启动
```
systemctl enable --now nginx mariadb
```

mysql初始化
```
mysql_secure_installation
```

登陆mysql创建数据库 配置完密码后修改配置文件
```
mysql -u root -p
CREATE DATABASE video CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
GRANT ALL PRIVILEGES ON video.* TO video@localhost IDENTIFIED BY 'vip#video123!';
FLUSH PRIVILEGES;
```

创建数据表
```
CREATE TABLE `visitor_logs` (
  `id` INT AUTO_INCREMENT PRIMARY KEY,
  `ip_address` VARCHAR(45) NOT NULL,
  `visit_time` DATETIME DEFAULT CURRENT_TIMESTAMP,
  `user_agent` TEXT
);
```

```
quit
```

下载代码
```
cd /var/www
wget -N https://github.com/hhttco/VipVideo/releases/download/v1.0.0/vv.zip -O ./vv.zip
unzip ./vv.zip -d /var/www/vv
chown -R www-data:www-data /var/www/vv && chmod -R 755 /var/www/vv
rm ./vv.zip
chmod +x /var/www/vv/vv
```

修改密码
```
vim config.json
```

配置守护进程
```
vim /etc/systemd/system/vv.service
```
```
[Unit]
Description=Vip video
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/var/www/vv
ExecStart=/var/www/sw/vv
Restart=always
MemoryLimit=50M

[Install]
WantedBy=multi-user.target
```


启动
```
systemctl daemon-reload
systemctl enable vv
systemctl start vv
systemctl status vv
```

配置域名 安装证书
```
vim /etc/nginx/conf.d/v.conf
```

```
server {
    server_name 域名;

    location / {
        # 转发请求到 Docker 映射出的本地端口
        proxy_pass http://127.0.0.1:8080;
        
        # 传递真实 IP（2026年标准：确保 Go 后端能拿到客户端真实 IP）
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

安装证书
```
apt -y install snapd && snap install core && snap refresh core && snap install --classic certbot && ln -s /snap/bin/certbot /usr/bin/certbot && certbot --nginx
```
```
systemctl reload nginx
```

##### 源代码编译
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
