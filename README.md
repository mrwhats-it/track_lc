# track_lc

To track leetcode qs solved amongst peers.

## how is works?

Github action schedules a script run at specified time duration(eg every 30 mins). The script pings lc api for total qs solved for each user and updates a json file containing usernames, questions solved and timestamp.

## wanna host it yourself?

1. fork the repo and create local clone
2. change usernames in scripts/update.go
3. push changes to remote fork
4. On your main repo page, go to Settings -> Actions -> General, under Workflow permissions tick "Read and write permissions" and save.
5. Go back to your main repo page, go to Actions tab and manually trigger workflow. Should show successful.
6. host frontend by connecting netlify/vercel with your fork along with the below settings

   | Setting               | Value                        |
   | --------------------- | ---------------------------- |
   | Runtime               | Not set                      |
   | Base directory        | `frontend`                   |
   | Package directory     | Not set                      |
   | Build command         | `npm run build`              |
   | Publish directory     | `frontend/dist`              |
   | Functions directory   | `frontend/netlify/functions` |
   | Deploy log visibility | Logs are public              |
   | Build status          | Active                       |

7. Enjoy

this project was inspired from [anga's lc tracker](https://battle.anga.codes/)
