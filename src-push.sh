eval "$(ssh-agent -s)"
ssh-add ~/.ssh/gitmota
git add .
git commit -m "Updated: `date +'%Y-%m-%d %H:%M:%S'`"
git push -f github main
git push -f gitlab main