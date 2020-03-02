<?php
use \Couchbase\ClusterOptions;
use \Couchbase\Cluster;
use \Couchbase\IncrementOptions;
use \Couchbase\DecrementOptions;

$opts = new ClusterOptions();
$opts->credentials("Administrator", "password");
$cluster = new Cluster("couchbase://192.168.1.101", $opts);

$bucket = $cluster->bucket("travel-sample");
$collection = $bucket->defaultCollection();
$binaryCollection = $collection->binary();


// Increment & Decrement are considered part of the 'binary' API
// and as such may still be a subject to change

// Create a document and assign it to 10 -- counter works atomically
// by first creating a document if it doesn't exist. If it exists,
// the same method will increment/decrement per the "delta" parameter
$key = "phpDevguideExampleCounter";
$opts = new IncrementOptions();
$opts->initial(10)->delta(2);

$res = $binaryCollection->increment($key, $opts);
// Should print 10
printf("Initialized Counter: %d\n", $res->content());

// Issue the same operation, increment value by 2 to 12
$res = $binaryCollection->increment($key, $opts);
// Should print 12
printf("Incremented Counter: %d\n", $res->content());

$opts = new DecrementOptions();
$opts->initial(10)->delta(4);
// Decrement value by 4 to 8
$res = $binaryCollection->decrement($key, $opts);
// Should print 8
printf("Decremented Counter: %d\n", $res->content());


// Output:
//
//     Initialized Counter: 10
//     Incremented Counter: 12
//     Decremented Counter: 8