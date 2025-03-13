#!/bin/bash

# 定义输出颜色
RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m' # 无颜色

# 检查root权限
if [ "$(id -u)" -ne 0 ]; then
    echo -e "${RED}错误：需要root权限执行此脚本${NC}"
    exit 1
fi

# 收集硬件信息函数
collect_hardware_info() {
    # 1. 主板信息
    BASEBOARD=$(dmidecode -t baseboard 2>/dev/null | grep "Serial Number" | awk -F': ' '{print $2}' | sed 's/\s*//g')
    [ -z "$BASEBOARD" ] && BASEBOARD=$(cat /sys/class/dmi/id/product_uuid 2>/dev/null)

    # 2. CPU信息
    CPU_ID=$(dmidecode -t processor 2>/dev/null | grep "ID" | awk -F': ' '{print $2}' | head -1 | sed 's/\s*//g')
    [ -z "$CPU_ID" ] && CPU_ID=$(grep -m1 "model name" /proc/cpuinfo | sha256sum | awk '{print $1}')

    # 3. 磁盘信息（获取第一个非虚拟磁盘的UUID）
    DISK_UUID=$(lsblk -o UUID,MOUNTPOINT -d -n -l 2>/dev/null | grep -w "/" | awk '{print $1}' | head -1)
    [ -z "$DISK_UUID" ] && DISK_UUID=$(blkid -s UUID -o value $(lsblk -o MOUNTPOINT,PKNAME -n -l | grep -w "/" | awk '{print $2}') 2>/dev/null)

    # 4. MAC地址（获取第一个物理网卡）
    NETWORK_MAC=$(ip link show | awk '/ether/ && !/lo/ {print $2; exit}' | tr -d ':')

    # 组合信息
    FINGERPRINT_RAW="${BASEBOARD:-NULL}|${CPU_ID:-NULL}|${DISK_UUID:-NULL}|${NETWORK_MAC:-NULL}"
    
    # 生成哈希指纹
    FINGERPRINT=$(echo -n "$FINGERPRINT_RAW" | sha256sum | awk '{print $1}')
    
    echo -e "${GREEN}硬件指纹生成成功${NC}"
    echo "指纹哈希：$FINGERPRINT"
    echo "原始信息：$FINGERPRINT_RAW"
}

# 主执行逻辑
echo "正在收集硬件信息..."
collect_hardware_info

# 可选保存到文件
read -p "是否保存到文件？[y/N] " SAVE
if [[ $SAVE =~ [Yy] ]]; then
    FILENAME="hardware_fingerprint_$(date +%s).txt"
    echo "指纹哈希：$FINGERPRINT" > $FILENAME
    echo "原始信息：$FINGERPRINT_RAW" >> $FILENAME
    echo -e "已保存至：${GREEN}$(pwd)/$FILENAME${NC}"
fi