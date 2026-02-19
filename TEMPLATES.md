StreamSpace Template Catalog
This document lists all available templates in the repository.

Statistics
Total Templates: 195
Total Categories: 50
Last Updated: 2025-11-18
Templates by Category
AI & ML (2 templates)
Faster-Whisper
Piper
Administration (2 templates)
Hishtory-Server
Openssh-Server
Audio & Video (15 templates)
Ardour
Audacity
Ffmpeg
Htpcmanager
Kdenlive
Minisatip
Openshot
Oscam
Shotcut
Synclounge
Tautulli
Tvheadend
Webgrabplus
Automation (2 templates)
Habridge
Homeassistant
Backup (4 templates)
Duplicati
Resilio-Sync
Rsnapshot
Syncthing
Communication (3 templates)
Limnoria
Ngircd
Znc
Design & Graphics (21 templates)
Bambustudio
Blender
Cura
Darktable
Digikam
FreeCAD
GIMP
Inkscape
KiCad
Krita
Lychee
Manyfold
Orcaslicer
Piwigo
Rawtherapee
Desktop Environments (3 templates)
Alpine Desktop (i3)
Ubuntu Desktop (KDE)
Ubuntu Desktop (XFCE)
Development (10 templates)
Code-Server
Github-Desktop
Gitqlient
Mariadb
Mysql-Workbench
Openvscode-Server
Phpmyadmin
Sqlitebrowser
Vscodium
Yaak
Docker (3 templates)
Modmanager
Qemu-Static
Socket-Proxy
File Management (9 templates)
Davos
Deluge
Filezilla
Flexget
Nzbget
Pyload-Ng
Qbittorrent
Sabnzbd
Transmission
File Sharing (4 templates)
Pairdrop
Projectsend
Pydio-Cells
Xbackbone
Gaming (13 templates)
Dogwalk
Dolphin
Duckstation
Flycast
Gzdoom
Luanti
Mame
Pcsx2
Retroarch
Rpcs3
Scummvm
Steamos
Xemu
Indexers (3 templates)
Jackett
Nzbhydra2
Prowlarr
Lifestyle (2 templates)
Babybuddy
Grocy
Media (14 templates)
Bazarr
Beets
Doplarr
Kometa
Lollypop
Mediaelch
Medusa
Ombi
Overseerr
Radarr
Sickgear
Sonarr
Spotube
Your_Spotify
Monitoring (6 templates)
Apprise-Api
Healthchecks
Librespeed
Smokeping
Speedtest-Tracker
Syslog-Ng
Network (2 templates)
Unifi-Network-Application
Wireshark
Productivity (22 templates)
Bitcoin-Knots
Bookstack
Budge
Calibre
Calibre-Web
Calligra
Cops
Dokuwiki
Grav
Hedgedoc
Joplin
Kavita
Lazylibrarian
Libreoffice
Monica
Obsidian
Onlyoffice
Raneto
Ubooquity
Wikijs
Wps-Office
Zotero
Remote Desktop (3 templates)
Remmina
Rustdesk
Webtop
Reverse Proxy (2 templates)
Nginx
Swag
Science & Education (2 templates)
Boinc
Foldingathome
Security (1 template)
Keepassxc
Web Browsers (14 templates)
Brave
Chrome
Chromium
Firefox
Librewolf
Msedge
Opera
Ungoogled-Chromium
Vivaldi
Zen
Image Sources
All templates use LinuxServer.io container images from lscr.io/linuxserver/:

Base Registry: lscr.io/linuxserver/
Icon Source: LinuxServer.io Docker Templates Repository
KasmVNC: Included for GUI applications (port 3000)
Web Apps: HTTP access (port 8080)
Usage
Deploy All Templates
kubectl apply -f https://raw.githubusercontent.com/JoshuaAFerguson/streamspace-templates/main/catalog.yaml
Deploy by Category
# Web Browsers
kubectl apply -f https://raw.githubusercontent.com/JoshuaAFerguson/streamspace-templates/main/web-browsers/

# Development Tools
kubectl apply -f https://raw.githubusercontent.com/JoshuaAFerguson/streamspace-templates/main/development/
Deploy Individual Template
kubectl apply -f https://raw.githubusercontent.com/JoshuaAFerguson/streamspace-templates/main/web-browsers/firefox.yaml
Contributing
See CONTRIBUTING.md for guidelines on adding new templates.

Links
Main Project: StreamSpace
LinuxServer.io: https://www.linuxserver.io/
Template Repository: streamspace-templates