-- @up

create keyspace "release"
  with replication = {'class': '$${placement_strategy}', $${strategy_options}}
  and durable_writes = true;

-- @down

drop keyspace "release";
