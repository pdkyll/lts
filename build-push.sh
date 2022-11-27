#!bin/bash
eval "$(ssh-agent -s)"
ssh-add ~/.ssh/gitmota
# синхронизируем исходники в gitlab и github
git add .
git commit -m "Updated: `date +'%Y-%m-%d %H:%M:%S'`"
git push -f github main
git push -f gitlab main
# создаем бинарные файлы
make build
# синхронизируем бинарники в gitlab и github
cd build
git add .
git commit -m "Updated: `date +'%Y-%m-%d %H:%M:%S'`"
#git remote add gitlab-web git@gitlab.com:m0ta/benefy-web.git
#git remote add github-web git@github.com:m0taru/benefy-web.git
#git branch -M main
git push -uf gitlab-build main
#git push -uf github-api main