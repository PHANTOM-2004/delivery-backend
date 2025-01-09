#!/bin/sh

# 遍历当前目录下的文件和子目录
for file in ./service/*; do
  env_file="$file/.env.example"
  if [ -e "$env_file" ]; then
    echo "find service .env.example in: $file"
    if [ ! -e "$file/.env" ]; then
      echo "copying $env_file"
      cp "$file/.env.example" "$file/.env"
    else
      echo ".env already exists, not copyint"
    fi
  fi
done
