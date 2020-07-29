-- select ic.item_id, ic.stop_time from item_children ic where item_id = $1;
select i.item_id,
       i.title,
       i.category_id,
       i.price,
       i.start_time,
       i.stop_time,
       ic.item_id   children_item_id,
       ic.stop_time children_stop_time
from item i
         left join item_children ic on i.item_id = ic.item_id
where i.item_id = $1;