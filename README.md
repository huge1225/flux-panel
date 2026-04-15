# flux-panel

本仓库为 `flux-panel` 的二次开发版本，安装资源与镜像已固定到 `huge1225` 命名空间。

## 一键安装（Docker Compose）

```bash
curl -L "https://raw.githubusercontent.com/huge1225/flux-panel/refs/heads/main/panel_install.sh" -o panel_install.sh && chmod +x panel_install.sh && ./panel_install.sh
```

## 默认端口

- 前端默认端口：`6366`
- 后端默认端口：`6365`

## 默认管理员

- 账号：`admin_user`
- 密码：`admin_user`

> 首次登录后请立即修改默认密码。

## 常用菜单功能

- 安装面板
- 更新面板
- 卸载面板
- 导出数据库备份
- 安装并配置反向代理（Caddy）

## 说明

- `panel_install.sh` 中安装地址已固定，不支持环境变量覆盖。
- `install.sh` 中 gost 下载地址已固定，不支持环境变量覆盖。
- `docker-compose-v4.yml` / `docker-compose-v6.yml` 中镜像地址已固定为：
  - `huge1225/springboot-backend:latest`
  - `huge1225/vite-frontend:latest`
