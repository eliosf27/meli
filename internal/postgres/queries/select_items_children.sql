select ic.item_id, ic.stop_time from item_children ic where item_id = $1;
