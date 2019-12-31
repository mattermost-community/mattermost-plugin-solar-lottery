```
/slottery user qualify -k webapp -l beginner -u @aaron.medina,@aaron.peterson,@aaron.ward,@albert.torres,@alice.johnston
/slottery user qualify -k server -l beginner -u @benjamin.bennett,@betty.campbell,@brenda.boyd,@christina.wilson,@craig.reed
/slottery user qualify -k lead -l beginner -u @emily.meyer,@eugene.rodriguez,@frances.elliott,@helen.hunter,@janice.armstrong

/slottery add rotation -r simple --period m -s 2019-12-17 --grace 1 --size 3
/slottery need -r simple -k webapp -l beginner --min 1 
/slottery need -r simple -k server -l beginner --min 1
/slottery need -r simple -k lead -l beginner --min 1 --max 1

/slottery join -r simple -u @aaron.medina,@aaron.peterson,@aaron.ward,@albert.torres,@alice.johnston,@benjamin.bennett,@betty.campbell,@brenda.boyd,@christina.wilson,@craig.reed,@emily.meyer,@eugene.rodriguez,@frances.elliott,@helen.hunter,@janice.armstrong

/slottery forecast guess -r simple -s 0 -n 6 --autofill
/slottery forecast rotation -r simple -s 0 -n 3 --sample 100




/slottery skill list

/slottery user qualify -k webapp -l beginner -u @aaron.peterson,@aaron.ward,@albert.torres,@alice.johnston
/slottery user qualify -k webapp -l intermediate -u @benjamin.bennett,@betty.campbell,@brenda.boyd,@christina.wilson,@craig.reed
/slottery user qualify -k server -l beginner -u @emily.meyer,@eugene.rodriguez,@frances.elliott,@helen.hunter,@jack.wheeler
/slottery user qualify -k server -l intermediate -u @jeremy.williamson,@jerry.ramos 
/slottery user qualify -k server -l advanced -u @benjamin.bennett,@janice.armstrong,@sysadmin
/slottery user qualify -k webapp -l advanced -u @aaron.medina,@janice.armstrong,@sysadmin
/slottery user qualify -k lead -l intermediate -u @aaron.medina,@benjamin.bennett,@emily.meyer,@janice.armstrong,@albert.torres,@alice.johnston

/slottery user show --users @aaron.medina,@albert.torres,@sysadmin

/slottery add rotation --rotation DEMO --period m --start 2019-12-17 --padding 1 --size 4
/slottery rotation need --rotation DEMO --skill webapp --level beginner --min 2 
/slottery rotation need --rotation DEMO --skill webapp --level intermediate --min 1 --max 1
/slottery rotation need --rotation DEMO --skill server --level beginner --min 2 --max 3
/slottery rotation need --rotation DEMO --skill server --level intermediate --min 1
/slottery rotation need --rotation DEMO --skill lead --level beginner --min 1 --max 1
/slottery rotation need --rotation DEMO --delete-need --skill lead --level beginner
/slottery rotation need --rotation DEMO --skill lead --level beginner --min 1 --max 1

/slottery join --rotation DEMO --users @aaron.medina,@aaron.peterson,@aaron.ward,@albert.torres,@alice.johnston,@benjamin.bennett,@betty.campbell,@brenda.boyd,@christina.wilson,@craig.reed,@emily.meyer,@eugene.rodriguez,@frances.elliott,@helen.hunter,@jack.wheeler,@janice.armstrong,@jeremy.williamson,@jerry.ramos,@sysadmin

/slottery show rotation --rotation DEMO

/slottery forecast schedule --rotation DEMO --start 2019-12-20 --shifts 3 --autofill
/slottery rotation update --rotation DEMO --size 5
/slottery rotation forecast DEMO --start 2019-12-20 --shifts 3 --autofill

/slottery add shift  --rotation DEMO --number 2 
/slottery shift join --rotation DEMO --number 2 --users @jack.wheeler,@janice.armstrong 

/slottery rotation archive DEMO
```