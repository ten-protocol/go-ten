# Obscuro Doc Site Staging Area

This is the staging area for the Obscuro Doc Site and it looks like [this](https://docsstage.obscu.ro/). This staging area is hosted in GavT's repo because of the way GitHub Pages, which is used to publish Obscuro docs, works.

## Adding New Doc Site Pages

[Jekyll](https://jekyllrb.com/) is used as a content deployment tool and Jekyll uses a Theme called [Minimal Mistakes](https://github.com/mmistakes/minimal-mistakes) to provide a simple UI. You will need to install Jekyll, but not Minimal Mistakes which is already configured in the repo.

1. [Install Jekyll](https://jekyllrb.com/docs/installation/)
2. Clone this repository: https://github.com/GavT/obsdocs.github.io
3. Create your new content as a markdown file in the `/docs` folder of the repo. Take care with the folder structure. As a general rule, new titles in the left hand navigation menu should have their content contained in a seperate subfolder under docs, for example, `/docs/testnet` contains all the markdown files relation to the testnet docs.
4. To have this new content shown in the left hand navigation menu you need to modify the file `/docs/_data/navigation.yml`. Follow the same format to add new headings and content titles. Remember to specify the file type as `.html` for your new markdown files, not `.md` when providing the URL.
5. Once you are happy with your markdown content run the jekyll build command to build the site pages:
  ```
  cd <cloned repo location>/docs
  bundle exec jekyll build
  ```
  This will create new html pages from your markdown pages.

6. Push your changes to GavT/obsdocs.github.io
7. Browse to https://docsstage.obscu.ro/ and check your content. Remember your browser will cache some of the pages so hit refresh a few times if it looks like the content is missing or the navigation menu is incorrect.

## Updating Content

This is the same as Adding New Doc Site Pages above but you can omit step 4 if you are not adding pages.
