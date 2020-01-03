```sh
# ABC-FS1
#       @aaron.medina - lead(3), server(3), webapp(3)
#       @aaron.peterson - server(1), webapp(2)
#       @aaron.ward - webapp(2)
#       @albert.torres - server(1)
#       @alice.johnston - server(2), webapp(1)
#
# DEF-PERF
#       @deborah.freeman - lead(2), server(3), webapp(2)
#       @diana.wells - sever(3)
#       @douglas.daniels - server(4)
#       @emily.meyer - server(3), webapp(1)
#       @eugene.rodriguez - server(2), webapp(1)
#       @frances.elliott - webapp(3)
#
# GHIJ-FS2
#       @helen.hunter - lead(3),webapp(1),server(1)
#       @janice.armstrong - server(2), webapp(1)
#       @jeremy.williamson - webapp(2)
#       @jerry.ramos - server(2), webapp(2)
#       @johnny.hansen - webapp(1), server(1)
#       @jonathan.watson - server(2)
#
# KL-SRE
#       @karen.austin - lead(2), sre(3)
#       @karen.martin - webapp(1), sre(1), build(3)
#       @kathryn.mills - server(2), sre(3)
#       @laura.wagner - sre(2)
#
# MN-INT
#       @margaret.morgan - lead(3), server(2)
#       @mark.rodriguez webapp(1), server(1)
#       @matthew.mendoza webapp(3), server(2)
#       @mildred.barnes webapp(2), server(2)
#       @nancy.roberts webapp(2), server(3)

#
# Open a window as @margaret.morgan

/lotto user qualify -k ABC -l intermediate -u @aaron.medina,@aaron.peterson,@aaron.ward,@albert.torres,@alice.johnston
/lotto user qualify -k DEF -l intermediate -u @deborah.freeman,@diana.wells,@douglas.daniels,@emily.meyer,@eugene.rodriguez,@frances.elliott
/lotto user qualify -k GHIJ -l intermediate -u @helen.hunter,@janice.armstrong,@jeremy.williamson,@jerry.ramos,@johnny.hansen,@jonathan.watson
/lotto user qualify -k KL -l intermediate -u @karen.austin,@karen.martin,@kathryn.mills,@laura.wagner
/lotto user qualify -k MN -l intermediate -u @margaret.morgan,@mark.rodriguez,@matthew.mendoza,@mildred.barnes,@nancy.roberts

/lotto user qualify -k lead -l intermediate -u @deborah.freeman,@karen.austin
/lotto user qualify -k lead -l advanced -u @aaron.medina,@helen.hunter,@margaret.morgan

/lotto user qualify -k server -l beginner -u @aaron.peterson,@albert.torres,@helen.hunter,@johnny.hansen,@mark.rodriguez
/lotto user qualify -k server -l intermediate -u @alice.johnston,@eugene.rodriguez,@janice.armstrong,@jerry.ramos,@jonathan.watson,@kathryn.mills,@margaret.morgan,@matthew.mendoza,@mildred.barnes
/lotto user qualify -k server -l advanced -u @aaron.medina,@deborah.freeman,@diana.wells,@emily.meyer,@nancy.roberts
/lotto user qualify -k server -l expert -u @douglas.daniels

/lotto user qualify -k webapp -l beginner -u @emily.meyer,@eugene.rodriguez,@helen.hunter,@janice.armstrong,@johnny.hansen,@karen.martin,@mark.rodriguez,@mildred.barnes
/lotto user qualify -k webapp -l intermediate -u @aaron.peterson,@aaron.ward,@deborah.freeman,@jeremy.williamson,@jerry.ramos,@nancy.roberts
/lotto user qualify -k webapp -l advanced -u @aaron.medina,@frances.elliott,@matthew.mendoza

/lotto user qualify -k sre -l beginner -u @karen.martin
/lotto user qualify -k sre -l intermediate -u @laura.wagner
/lotto user qualify -k sre -l advanced -u @karen.austin,@kathryn.mills

/lotto user qualify -k build -l advanced -u @karen.martin

/lotto rotation add -r icebreaker --period w -s 2019-12-17 --grace 3 --size 2

/lotto rotation leave -r icebreaker -u @aaron.medina,@aaron.peterson,@aaron.ward,@albert.torres,@alice.johnston,@deborah.freeman,@diana.wells,@douglas.daniels,@emily.meyer,@eugene.rodriguez,@frances.elliott,@helen.hunter,@janice.armstrong,@jeremy.williamson,@jerry.ramos,@johnny.hansen,@jonathan.watson,@karen.austin,@karen.martin,@kathryn.mills,@laura.wagner,@margaret.morgan,@mark.rodriguez,@matthew.mendoza,@mildred.barnes,@nancy.roberts

/lotto rotation join -r icebreaker -s 2019-12-17 -u @aaron.medina,@aaron.peterson,@aaron.ward,@albert.torres,@alice.johnston,@deborah.freeman,@diana.wells,@douglas.daniels,@emily.meyer,@eugene.rodriguez,@frances.elliott,@helen.hunter,@janice.armstrong,@jeremy.williamson,@jerry.ramos,@johnny.hansen,@jonathan.watson,@karen.austin,@karen.martin,@kathryn.mills,@laura.wagner,@margaret.morgan,@mark.rodriguez,@matthew.mendoza,@mildred.barnes,@nancy.roberts

/lotto rotation guess -r icebreaker -s 0 -n 3 --autofill
/lotto rotation forecast -r icebreaker -s 0 -n 20 --sample 100
/lotto rotation autopilot -r icebreaker --notify 3 --fill-before 5
/lotto rotation autopilot -r icebreaker --debug-run 2019-12-11






/lotto skill list

/lotto user qualify -k webapp -l beginner -u @aaron.peterson,@aaron.ward,@albert.torres,@alice.johnston
/lotto user qualify -k webapp -l intermediate -u @benjamin.bennett,@betty.campbell,@brenda.boyd,@christina.wilson,@craig.reed
/lotto user qualify -k server -l beginner -u @emily.meyer,@eugene.rodriguez,@frances.elliott,@helen.hunter,@jack.wheeler
/lotto user qualify -k server -l intermediate -u @jeremy.williamson,@jerry.ramos 
/lotto user qualify -k server -l advanced -u @benjamin.bennett,@janice.armstrong,@sysadmin
/lotto user qualify -k webapp -l advanced -u @aaron.medina,@janice.armstrong,@sysadmin
/lotto user qualify -k lead -l intermediate -u @aaron.medina,@benjamin.bennett,@emily.meyer,@janice.armstrong,@albert.torres,@alice.johnston

/lotto user show --users @aaron.medina,@albert.torres,@sysadmin

/lotto add rotation --rotation DEMO --period m --start 2019-12-17 --padding 1 --size 4
/lotto rotation need --rotation DEMO --skill webapp --level beginner --min 2 
/lotto rotation need --rotation DEMO --skill webapp --level intermediate --min 1 --max 1
/lotto rotation need --rotation DEMO --skill server --level beginner --min 2 --max 3
/lotto rotation need --rotation DEMO --skill server --level intermediate --min 1
/lotto rotation need --rotation DEMO --skill lead --level beginner --min 1 --max 1
/lotto rotation need --rotation DEMO --delete-need --skill lead --level beginner
/lotto rotation need --rotation DEMO --skill lead --level beginner --min 1 --max 1

/lotto join --rotation DEMO --users @aaron.medina,@aaron.peterson,@aaron.ward,@albert.torres,@alice.johnston,@benjamin.bennett,@betty.campbell,@brenda.boyd,@christina.wilson,@craig.reed,@emily.meyer,@eugene.rodriguez,@frances.elliott,@helen.hunter,@jack.wheeler,@janice.armstrong,@jeremy.williamson,@jerry.ramos,@sysadmin

/lotto show rotation --rotation DEMO

/lotto forecast schedule --rotation DEMO --start 2019-12-20 --shifts 3 --autofill
/lotto rotation update --rotation DEMO --size 5
/lotto rotation forecast DEMO --start 2019-12-20 --shifts 3 --autofill

/lotto add shift  --rotation DEMO --number 2 
/lotto shift join --rotation DEMO --number 2 --users @jack.wheeler,@janice.armstrong 

/lotto rotation archive DEMO
```