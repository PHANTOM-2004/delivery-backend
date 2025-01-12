sudo redis-server ./redis_master.conf
sudo redis-server ./redis_slave.conf
sudo redis-sentinel ./sentinel.conf --sentinel
