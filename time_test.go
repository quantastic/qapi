package qapi

//func Test_UTCTime(t *testing.T) {
//tests := []struct {
//Input   string
//Want    string
//WantErr bool
//}{
//{
//Input: `"2014-08-11T10:23:51Z"`,
//Want:  `"2014-08-11T10:23:51Z"`,
//},
//{
//Input: `"2014-08-11T11:23:51+01:00"`,
//Want:  `"2014-08-11T10:23:51Z"`,
//},
//{
//Input: `"2014-08-11T10:23:51.3245Z"`,
//Want:  `"2014-08-11T10:23:51Z"`,
//},
//}
//for i, test := range tests {
//ut := UTCTime{}
//if err := json.Unmarshal([]byte(test.Input), &ut); test.WantErr {
//if err == nil {
//t.Errorf("test %d: Wanted error, but got none.", i)
//}
//continue
//} else if err != nil {
//t.Errorf("test %d: Unmarshal: %s", i, err)
//continue
//}
//got, err := json.Marshal(ut)
//if err != nil {
//t.Errorf("test %d: Marshal: %s", i, err)
//continue
//}
//gotStr := string(got)
//if diff := pretty.Compare(gotStr, test.Want); diff != "" {
//t.Errorf("test %d:\n%s", i, diff)
//continue
//}
//}
//}
