package universe_test


import "testing"
import "csv"

option now = () => 2030-01-01T00:00:00Z

input =
    "
#datatype,string,long,dateTime:RFC3339,string,string,string,double
#group,false,false,false,true,true,true,false
#default,_result,,,,,,
,result,table,_time,_measurement,host,_field,_value
,,0,2018-05-22T19:53:26Z,system,host.local,load1,1.83
,,0,2018-05-22T19:53:36Z,system,host.local,load1,1.63
,,1,2018-05-22T19:53:26Z,system,host.local,load3,1.72
,,2,2018-05-22T19:53:26Z,system,host.local,load4,1.77
,,2,2018-05-22T19:53:36Z,system,host.local,load4,1.78
,,2,2018-05-22T19:53:46Z,system,host.local,load4,1.77
"
output =
    "
#datatype,string,long,dateTime:RFC3339,dateTime:RFC3339,dateTime:RFC3339,string,string,string,double
#group,false,false,true,true,false,true,true,true,false
#default,_result,,,,,,,,
,result,table,_start,_stop,_time,_measurement,host,_field,_value
,,0,2018-05-22T19:53:26Z,2030-01-01T00:00:00Z,2018-05-22T19:53:26Z,system,host.local,load4,1.77
,,0,2018-05-22T19:53:26Z,2030-01-01T00:00:00Z,2018-05-22T19:53:46Z,system,host.local,load4,1.77
"

testcase merge_filter_evaluate {
    got =
        csv.from(csv: input)
            |> testing.load()
            |> range(start: 2018-05-22T19:53:26Z)
            |> filter(fn: (r) => r["_value"] == 1.77)
            |> filter(fn: (r) => r["_field"] == "load4")
    want = csv.from(csv: output)

    testing.diff(got, want)
}
