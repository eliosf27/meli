select i.item_id, i.title, i.category_id, i.price, i.start_time, i.stop_time
from item i
where item_id = $1;