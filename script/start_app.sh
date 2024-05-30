#!/bin/bash

# 需要手动修改
appNmae="github.com/yanshicheng/ikubeops-gin-demo"
newAppShortName="github.com/yanshicheng"
newAppName="${appNmae}/apps"
# shellcheck disable=SC2034
NewRegistryUrl="harbor.ikubeops.local/public"

# shellcheck disable=SC2034
OldRegistryUrl="harbor.ikubeops.local/public"
# 完整的项目名字
oldProjectFullNmae="github.com/yanshicheng/ikubeops-gin-demo"
oldProjectName="ikubeops-gin-demo"
oldProjectShortName="github.com/yanshicheng"

# 获取脚本的绝对路径
scriptDir=$(cd "$(dirname "$0")"; pwd)
# 项目的根目录
baseDir=$(cd $scriptDir/../; pwd)
# 应用目录
appDir=$baseDir/apps
# 示例应用目录
demoDir=$baseDir/template/apps/demo

# 初始化 packageName 变量
packageName=""
oldAppName="${appNmae}/template/apps/demo"
# shellcheck disable=SC2034
oldName="demo"

projectName=""
# 显示帮助信息
show_help() {
    echo "使用方法: $0 [选项]"
    echo ""
    echo "选项:"
    echo "  -h               显示此帮助信息"
    echo "  -i projectName   初始化项目, 请先修改 AppNamme 值"
    echo "  -a packageName   指定要初始化的包名"
    exit 0
}

check_package_name() {
# 检查是否提供了必需的 -a 选项
  if [ -z "$packageName" ]; then
      echo "错误：请使用 -a 选项并提供 packageName。"
      exit 1
  fi
}

# 检查路径函数
check_paths() {
    echo "检查路径："
    if [ -d "$appDir" ]; then
        echo "appDir 存在: $appDir"
    else
        echo "appDir 请确定确的 appDir 路径..."
        exit 1
    fi
    if [ -d "$demoDir" ]; then
        echo "demoDir 存在: $demoDir"
    else
        echo "demoDir 请确定确的 demoDir 路径: $demoDir"
        exit 1
    fi
    if [ -d "$appDir/$packageName" ]; then
        echo "错误：目录 $appDir/$packageName 已存在。"
        exit 1
    fi
}

if [ "$#" -eq 0 ]; then
    echo "错误：请使用 -h 选项获取帮助信息。"
    exit 1
fi


copy_dir() {
  echo "创建app: $appDir/$packageName"
  cp -r $demoDir $appDir/$packageName
  if [ $? -eq 0 ]; then
      echo "创建成功！"
  else
      echo "创建失败！"
      exit 1
  fi
}

inject_app() {
  echo "注入app: $appDir/$packageName"
  # 替换 packageName
  sed -i  "s/$oldName/$packageName/g" "$appDir/$packageName/app.go"
  sed -i  "s/$oldName.AppName/$packageName.AppName/g" "$appDir/$packageName/logic/logic.go"
  sed -i  "s/$oldName.AppName/$packageName.AppName/g" "$appDir/$packageName/handler/router.go"
  # 替换文件名
  find "$appDir/$packageName/handler" -depth -name "*$oldName*" -exec bash -c '
  for f; do
      newname=$(echo "$f" | sed "s/'"$oldName"'/'"$packageName"'/g")
      # 执行重命名操作
      if [ "$f" != "$newname" ]; then
          mv "$f" "$newname"
      fi
  done
  ' bash {} +
  # 替换文件名
  find "$appDir/$packageName/logic" -depth -name  "$oldName" -exec bash -c '
  for f; do
      newname=$(echo "$f" | sed "s/'"$oldName"'/'"$packageName"'/g")
      # 执行重命名操作
      if [ "$f" != "$newname" ]; then
          mv "$f" "$newname"
      fi
  done
  ' bash {} +
  # 更换所有 import name
  sNmae=$newAppName/$packageName
  find "$appDir/$packageName/" -type f -exec sed -i  "s@$oldAppName@$sNmae@g" {} \;
}

inject_app_mac() {
  echo "注入app: $appDir/$packageName"
  # 替换 packageName
  sed -i '' "s/$oldName/$packageName/g" "$appDir/$packageName/app.go"
  sed -i '' "s/$oldName.AppName/$packageName.AppName/g" "$appDir/$packageName/logic/logic.go"
  sed -i '' "s/$oldName.AppName/$packageName.AppName/g" "$appDir/$packageName/handler/router.go"
  find "$appDir/$packageName/handler" -depth -name "*$oldName*" -exec bash -c '
  for f; do
      echo $f
      newname=$(echo "$f" | sed "s/'"$oldName"'/'"$packageName"'/g")
      # 执行重命名操作
      if [ "$f" != "$newname" ]; then
          mv "$f" "$newname"
      fi
  done
  ' bash {} +
  find "$appDir/$packageName/logic" -depth -name "*$oldName*" -exec bash -c '
  for f; do
      newname=$(echo "$f" | sed "s/'"$oldName"'/'"$packageName"'/g")
      # 执行重命名操作
      if [ "$f" != "$newname" ]; then
          mv "$f" "$newname"
      fi
  done
  ' bash {} +
  sNmae=$newAppName/$packageName
  find "$appDir/$packageName/" -type f -exec sed -i '' "s@$oldAppName@$sNmae@g" {} \;
}

inject_apps_logic() {
  echo "注入全局 logic"
  msg1="    _ \"${appNmae}/apps/${packageName}/handler\""
  msg2="    _ \"${appNmae}/apps/${packageName}/logic\""
  msg3="    _ \"${appNmae}/apps/${packageName}/models\""
awk -v msg1="$msg1" -v msg2="$msg2" -v msg3="$msg3" '
{
  if (/\)/) {
    last=NR;
    last_line=$0;
  }
  hold[NR]=$0;
}
END {
  for (i=1; i<last; i++) {
    print hold[i];
  }
  sub(/\)/, "\n" msg1 "\n" msg2 "\n" msg3 "\n&", last_line);
  print last_line;
  for (i=last+1; i<=NR; i++) {
    print hold[i];
  }
}' "${appDir}/all/logic.go" > "${appDir}/all/temp_logic.go" && mv "${appDir}/all/temp_logic.go" "${appDir}/all/logic.go"
}

start_app() {
  sys=$(uname)
  check_paths
  copy_dir

  if [ "$sys" = "Darwin" ]; then
      inject_app_mac
  else
      inject_app
  fi

  inject_apps_logic
  echo "新增 $packageName App 成功!!!"
}

change_project() {
  echo "Linux 修改项目配置"
  find "$baseDir/apps/" -type f -exec sed -i  "s@$oldProjectFullNmae@$appNmae@g" {} \;
  find "$baseDir/cmd/" -type f -exec sed -i  "s@$oldProjectFullNmae@$appNmae@g" {} \;
  find "$baseDir/common/" -type f -exec sed -i  "s@$oldProjectFullNmae@$appNmae@g" {} \;
  find "$baseDir/config/" -type f -exec sed -i  "s@$oldProjectFullNmae@$appNmae@g" {} \;
  find "$baseDir/global/" -type f -exec sed -i  "s@$oldProjectFullNmae@$appNmae@g" {} \;
  find "$baseDir/template/" -type f -exec sed -i  "s@$oldProjectFullNmae@$appNmae@g" {} \;
  find "$baseDir/middleware/" -type f -exec sed -i  "s@$oldProjectFullNmae@$appNmae@g" {} \;
  find "$baseDir/router/" -type f -exec sed -i  "s@$oldProjectFullNmae@$appNmae@g" {} \;
  find "$baseDir/settings/" -type f -exec sed -i  "s@$oldProjectFullNmae@$appNmae@g" {} \;
  find "$baseDir/test/" -type f -exec sed -i  "s@$oldProjectFullNmae@$appNmae@g" {} \;
  find "$baseDir/utils/" -type f -exec sed -i  "s@$oldProjectFullNmae@$appNmae@g" {} \;
  find "$baseDir/version/" -type f -exec sed -i  "s@$oldProjectFullNmae@$appNmae@g" {} \;
  sed -i  "s@$oldProjectName@$projectName@g" "$baseDir/version/version.go"
  sed -i  "s@$oldProjectFullNmae@$appNmae@g" "$baseDir/main.go"
  sed -i  "s@$oldProjectShortName@$newAppShortName@g" "$baseDir/Makefile"
  # 修改 makefile
  sed -i  "s@$oldProjectName@$projectName@g" "$baseDir/Makefile"
  sed -i  "s@$OldRegistryUrl@$NewRegistryUrl@g" "$baseDir/Makefile"
  sed -i  "s@$oldProjectFullNmae@$appNmae@g" "$baseDir/go.mod"
  sed -i  "s@$oldProjectName@$projectName@g" "$baseDir/Dockerfile"
}


change_project_mac() {
  echo "Mac 修改项目配置"
  find "$baseDir/apps/" -type f -exec sed -i '' "s@$oldProjectFullNmae@$appNmae@g" {} \;
  find "$baseDir/cmd/" -type f -exec sed -i '' "s@$oldProjectFullNmae@$appNmae@g" {} \;
  find "$baseDir/common/" -type f -exec sed -i '' "s@$oldProjectFullNmae@$appNmae@g" {} \;
  find "$baseDir/template/" -type f -exec sed -i '' "s@$oldProjectFullNmae@$appNmae@g" {} \;
  find "$baseDir/config/" -type f -exec sed -i '' "s@$oldProjectFullNmae@$appNmae@g" {} \;
  find "$baseDir/global/" -type f -exec sed -i '' "s@$oldProjectFullNmae@$appNmae@g" {} \;
  find "$baseDir/middleware/" -type f -exec sed -i '' "s@$oldProjectFullNmae@$appNmae@g" {} \;
  find "$baseDir/router/" -type f -exec sed -i '' "s@$oldProjectFullNmae@$appNmae@g" {} \;
  find "$baseDir/settings/" -type f -exec sed -i '' "s@$oldProjectFullNmae@$appNmae@g" {} \;
  find "$baseDir/test/" -type f -exec sed -i '' "s@$oldProjectFullNmae@$appNmae@g" {} \;
  find "$baseDir/utils/" -type f -exec sed -i '' "s@$oldProjectFullNmae@$appNmae@g" {} \;
  find "$baseDir/version/" -type f -exec sed -i '' "s@$oldProjectFullNmae@$appNmae@g" {} \;
  sed -i '' "s@$oldProjectFullNmae@$appNmae@g" "$baseDir/main.go"
  sed -i '' "s@$oldProjectName@$projectName@g" "$baseDir/version/version.go"
  sed -i '' "s@$oldProjectShortName@$newAppShortName@g" "$baseDir/Makefile"
  sed -i '' "s@$oldProjectName@$projectName@g" "$baseDir/Makefile"
  sed -i '' "s@$OldRegistryUrl@$NewRegistryUrl@g" "$baseDir/Makefile"
  sed -i '' "s@$oldProjectFullNmae@$appNmae@g" "$baseDir/go.mod"
  sed -i '' "s@$oldProjectName@$projectName@g" "$baseDir/Dockerfile"
}
change_git_url() {
  rm -rf $baseDir/.git
  git init  &>/dev/null
  git branch -M main &>/dev/null
  git add . &>/dev/null
  git commit -m "first commit" &>/dev/null
  git remote add origin  https://$appNmae
  git remote -v
  echo "上传代码请执行: git push -u origin main"
}

change_app_desc() {
  sys=`uname`
  if [ "$sys" = "Darwin" ]; then
        change_project_mac
  else
        change_project
  fi
  change_git_url
  # 判断 $?
  if [ $? -eq 0 ]; then
      echo "修改项目配置成功！"
  else
      echo "修改项目配置失败！"
      exit 1
  fi
}

OptString=":ha:i:"
while getopts $OptString opt; do
    case $opt in
        h)
            show_help
            exit 0
            ;;
        a)
            if [ -z "$OPTARG" ]; then
                echo "选项 -a 需要一个参数。" >&2
                exit 1
            fi
            echo "选项: $opt, 值: $OPTARG, 索引: $OPTIND"
            packageName=$OPTARG
            # 假设 check_package_name 和 start_app 是已定义的函数
            check_package_name
            start_app
            ;;
        i)
            if [ -z "$OPTARG" ]; then
                echo "选项 -i 需要一个参数。" >&2
                exit 1
            fi
            echo "选项: $opt, 值: $OPTARG, 索引: $OPTIND"
            projectName=$OPTARG
            # 假设 change_app_desc 是已定义的函数
            change_app_desc
            ;;
        \?)
            echo "无效选项: -$OPTARG。请使用 -h 获取帮助信息。" >&2
            exit 1
            ;;
        :)
            echo "选项 -$OPTARG 需要一个参数。" >&2
            exit 1
            ;;
    esac
done