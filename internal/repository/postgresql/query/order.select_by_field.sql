select 
    id, id_customer, date, status 
from orders 
where
    case 
        when $1::varchar is not null and char_length($1) != 0 THEN
            id = $1
        else 1=1
    end
and
    case 
        when $2::varchar is not null and char_length($2) != 0 THEN
            id_customer = $2
        else 1=1
    end
and
    case 
        when $3::varchar is not null and char_length($3) != 0 THEN
            status = $3
        else 1=1
    end
    