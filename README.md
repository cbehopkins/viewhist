# viewhist
View (deep) History of websites

Problem:
Many websites have rss feeds. Most RSS feeds only show the recent history. For feeds where you want to start at the beginning this can mean a lot of trouble navigating to the start of their feeds and *troublesome* archives to navigate through. Even worse you can miss posts as new ones are posted.

Initial Aim:
Build a webserver that fetches the history from remote websites then drip feeds it to you.
The initial implementation works with Tumblr which has one of the worst history sections - great if you want to see things in reverse time order, but terrible to go back to the beginning and work forwards

So we implement a tumblr app that fetches the specified posts, extracts (in the first case) the text and renders a webpage with a number of posts on it.
One can then refresh the webpage and get the next bunch of posts
You can change the number you get in each fetch.

Trying to keep this a simple app for personal use - so I'm REALLY not concerned with things like user management security as (at the moemnt) it is only visible within my internal network anyway.
There is a json file (not checked in) where you specify the user config and therefore what you want to pull and it does so on every view page hit.
Currently only working with one page that we are reading from tumblr, next job it to expand it so that I can catch up with multiple tumblrs at once
Then to add the ability to get from other websites - I have my eye set on a number of webcomics with a massive back history

This probably needs a change to the back end so that rather than fetch everything on the view page request it can pre-load the items for the next page inbetween requests. This has potential scalability problems for large services, but not for me doing it for myself
