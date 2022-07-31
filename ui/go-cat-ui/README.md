# UI for go-cat

This is a statically generated frontend for go-cat environment list 
which is built with Typescript and Svelte. 

## Recommended IDE Setup

[VS Code](https://code.visualstudio.com/) + [Svelte](https://marketplace.visualstudio.com/items?itemName=svelte.svelte-vscode).

## Usage 
On your CI, attach the following 

```bash 
git clone https://gitlab.com/sorcero/community/go-cat 

cd go-cat/ui/go-cat-ui 
cp $SOURCE/infra.json src/.
yarn install 
yarn run build 

# set PUBLIC_DIR to the place you would want to export it
cp -r dist $PUBLIC_DIR/.
```
