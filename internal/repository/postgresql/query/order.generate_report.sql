select o.id, c.name, o.date, sum(p.price), o.status from orders o 
    inner join product_order po on o.id = po.id_order
    inner join customer c on o.id_customer = c.id
    inner join product p on po.id_product = p.id
group by o.id, c.name, o.date, o.status