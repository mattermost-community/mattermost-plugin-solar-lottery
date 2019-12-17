```
/slottery skill list

/slottery skill add webapp --level beginner --users @aaron.peterson,@aaron.ward,@albert.torres,@alice.johnston
/slottery skill add webapp --level intermediate --users @benjamin.bennett,@betty.campbell,@brenda.boyd,@christina.wilson,@craig.reed
/slottery skill add server --level beginner --users @emily.meyer,@eugene.rodriguez,@frances.elliott,@helen.hunter,@jack.wheeler
/slottery skill add server --level intermediate --users @jeremy.williamson,@jerry.ramos 
/slottery skill add server --level advanced --users @benjamin.bennett,@janice.armstrong,@sysadmin
/slottery skill add webapp --level advanced --users @aaron.medina,@janice.armstrong,@sysadmin
/slottery skill add lead --level intermediate --users @aaron.medina,@benjamin.bennett,@emily.meyer,@janice.armstrong,@albert.torres,@alice.johnston

/slottery user show --users @aaron.medina,@albert.torres,@sysadmin

/slottery rotation create DEMO --period m --start 2019-12-17 --padding 1 --size 4
/slottery rotation need DEMO --need fe1 --skill webapp --level beginner --min 2 
/slottery rotation need DEMO --need fe2 --skill webapp --level intermediate --min 1 --max 1
/slottery rotation need DEMO --need server1 --skill server --level beginner --min 2 --max 3
/slottery rotation need DEMO --need server2 --skill server --level intermediate --min 1
/slottery rotation need DEMO --need lead --skill lead --level beginner --min 1 --max 1
/slottery rotation need DEMO --remove-need lead
/slottery rotation need DEMO --need lead --skill lead --level beginner --min 1 --max 1

/slottery join DEMO --users @aaron.medina,@aaron.peterson,@aaron.ward,@albert.torres,@alice.johnston,@benjamin.bennett,@betty.campbell,@brenda.boyd,@christina.wilson,@craig.reed,@emily.meyer,@eugene.rodriguez,@frances.elliott,@helen.hunter,@jack.wheeler,@janice.armstrong,@jeremy.williamson,@jerry.ramos,@sysadmin

/slottery rotation show DEMO

/slottery rotation forecast DEMO --start 2019-12-20 --shifts 3 --autofill
/slottery rotation update DEMO --size 5
/slottery rotation forecast DEMO --start 2019-12-20 --shifts 3 --autofill

/slottery rotation archive DEMO
```