# inventory
Inventory Management System

```
$ inventory add 10 sd-card 10 rpi3
added 10 sd-card
added 10 rpi3

$ inventory remove 1 rpi3
removed 1 rpi3

$ inventory report
ITEMS
sd-card: 10
rpi3:    9

$ inventory set basic 1 sd-card 1 rpi3
basic is 1 sd-card & 1 rpi3

$ inventory report
ITEMS
sd-card: 10
rpi3:    9

SETS
basic:   9

$ inventory remove 1 basic
removed 1 sd-card
removed 1 rpi3

$ inventory report
ITEMS
sd-card: 9
rpi3:    8

SETS
basic:   8

$ inventory unset basic
basic is unset

$ inventory report
ITEMS
sd-card: 9
rpi3:    8
```

# ToDO

- [x] inventory add
- [x] inventory report
- [ ] inventory remove
- [ ] inventory count (reconcile)
- [ ] inventory set
- [ ] inventory unset