package main

import (
	"4d63.com/tz"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func TestGetPlanning(t *testing.T) {
	loc, _ := tz.LoadLocation("Europe/Amsterdam")
	type args struct {
		schedule Schedule
		date     time.Time
	}
	tests := []struct {
		name    string
		status  int
		json    string
		args    args
		wantErr bool
		want    *Planning
	}{
		{
			name: "Planning for Schedule",
			args: args{
				schedule: Schedule{
					GroupName:   "GroupName",
					ParameterId: "P1",
					GroupId:     "G1",
				},
				date: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			status: http.StatusOK,
			json: `
			{
			  "enableManualPlanning": false,
			  "enablePrimarySchedule": true,
			  "predefinedTimeSlots": [
				{
				  "start": "2021-05-31T00:00:00Z",
				  "end": "2021-06-01T00:00:00Z",
				  "minMembers": 1,
				  "maxMembers": 2
				}
			  ],
			  "baseTimeSlots": [
				{
				  "members": [
					"a9f656bf-85af-415b-807c-81728f255f03"
				  ],
				  "start": "2021-05-31T00:00:00Z",
				  "end": "2021-06-01T00:00:00Z",
				  "minMembers": 1,
				  "maxMembers": 2
				}
			  ],
			  "primaryTimeSlots": []
			}
`,
			wantErr: false,
			want: &Planning{
				PrimaryTimeSlots: []Slot{},
				BaseTimeSlots: []Slot{
					{
						Start: time.Date(2021, 5, 31, 0, 0, 0, 0, loc),
						End:   time.Date(2021, 6, 1, 0, 0, 0, 0, loc),
						Members: []string{
							"a9f656bf-85af-415b-807c-81728f255f03",
						},
					},
				},
			},
		},
		{
			name: "Planning for Schedule returns error",
			args: args{
				schedule: Schedule{
					GroupName:   "GroupName",
					ParameterId: "P1",
					GroupId:     "G1",
				},
				date: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			status:  http.StatusInternalServerError,
			json:    "",
			wantErr: true,
			want:    nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.status)
				w.Write([]byte(tt.json))
			}))
			defer ts.Close()
			nerveCentreBaseUrl = ts.URL

			got, err := GetPlanning(tt.args.schedule, tt.args.date)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetPlanning() = %v, want %v", got, tt.want)
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("GetPlanning() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetSchedules(t *testing.T) {
	tests := []struct {
		name string
		want *[]Schedule
	}{
		{
			name: "Single Schedule",
			want: &[]Schedule{
				{
					GroupName:   "GroupName",
					ParameterId: "P1",
					GroupId:     "G1",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				body, _ := json.Marshal(tt.want)
				w.WriteHeader(http.StatusOK)
				w.Write(body)
			}))
			defer ts.Close()
			nerveCentreBaseUrl = ts.URL
			if got := GetSchedules(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSchedules() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetUsers(t *testing.T) {
	tests := []struct {
		name string
		want *[]User
	}{
		{
			"Single User",
			&[]User{
				{
					Id:        "1",
					FirstName: "alice",
					LastName:  "",
				},
			},
		},
		{
			"Two Users",
			&[]User{
				{
					Id:        "1",
					FirstName: "alice",
					LastName:  "",
				},
				{
					Id:        "2",
					FirstName: "bob",
					LastName:  "",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				body, _ := json.Marshal(tt.want)
				w.WriteHeader(http.StatusOK)
				w.Write(body)
			}))
			defer ts.Close()
			nerveCentreBaseUrl = ts.URL
			if got := GetUsers(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUsers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLogin(t *testing.T) {
	type args struct {
		username string
		password string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Success",
			args: args{
				username: "bob",
				password: "alice",
			},
			wantErr: false,
		},
		{
			name: "Failure",
			args: args{
				username: "bob",
				password: "alice",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.wantErr {
					w.WriteHeader(http.StatusForbidden)
				} else {
					w.Header().Set("Location", nerveCentreBaseUrl+"?ReturnUrl=~%2f&State=1234567890")
					w.WriteHeader(http.StatusFound)
				}
			}))
			defer ts.Close()
			nerveCentreBaseUrl = ts.URL

			if err := Login(tt.args.username, tt.args.password); (err != nil) != tt.wantErr {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPlanning_HasMembers(t *testing.T) {
	type fields struct {
		BaseTimeSlots    []Slot
		PrimaryTimeSlots []Slot
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "Single Slot with Members",
			fields: fields{
				BaseTimeSlots: []Slot{
					{
						Start: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
						End:   time.Date(2021, 1, 1, 23, 59, 59, 0, time.UTC),
						Members: []string{
							"1",
						},
					},
				},
			},
			want: true,
		},
		{
			name: "Double Slot with Members",
			fields: fields{
				BaseTimeSlots: []Slot{
					{
						Start: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
						End:   time.Date(2021, 1, 1, 23, 59, 59, 0, time.UTC),
						Members: []string{
							"1", "2",
						},
					},
					{
						Start: time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC),
						End:   time.Date(2021, 1, 2, 23, 59, 59, 0, time.UTC),
						Members: []string{
							"2", "3",
						},
					},
				},
			},
			want: true,
		},
		{
			name: "Single Slot without Members",
			fields: fields{
				BaseTimeSlots: []Slot{
					{
						Start:   time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
						End:     time.Date(2021, 1, 1, 23, 59, 59, 0, time.UTC),
						Members: []string{},
					},
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			planning := &Planning{
				BaseTimeSlots:    tt.fields.BaseTimeSlots,
				PrimaryTimeSlots: tt.fields.PrimaryTimeSlots,
			}
			if got := planning.HasMembers(); got != tt.want {
				t.Errorf("HasMembers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSlot_GetMembers(t *testing.T) {
	type args struct {
		users *[]User
	}
	users := &[]User{
		{
			Id:        "1",
			FirstName: "Alice",
			LastName:  "",
		},
		{
			Id:        "2",
			FirstName: "Bob",
			LastName:  "",
		},
		{
			Id:        "3",
			FirstName: "Clare",
			LastName:  "",
		},
	}
	tests := []struct {
		name   string
		fields Slot
		args   args
		want   []string
	}{
		{
			name: "Single Slot",
			args: args{
				users: users,
			},
			fields: Slot{
				Start: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2021, 1, 1, 23, 59, 59, 0, time.UTC),
				Members: []string{
					"1",
				},
			},
			want: []string{"Alice"},
		},
		{
			name: "Double Slot",
			args: args{
				users: users,
			},
			fields: Slot{
				Start: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2021, 1, 1, 23, 59, 59, 0, time.UTC),
				Members: []string{
					"1", "2",
				},
			},
			want: []string{"Alice", "Bob"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fields.GetMembers(tt.args.users); !sameStringSlice(got, tt.want) {
				t.Errorf("GetMembers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPlanning_GetActiveSlot(t *testing.T) {
	type fields struct {
		BaseTimeSlots    []Slot
		PrimaryTimeSlots []Slot
	}
	type args struct {
		time time.Time
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Slot
	}{
		{
			name: "Outside of Slot",
			args: args{
				time: time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC),
			},
			fields: fields{
				BaseTimeSlots: []Slot{
					{
						Start: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
						End:   time.Date(2021, 1, 1, 23, 59, 59, 0, time.UTC),
						Members: []string{
							"1",
						},
					},
				},
			},
			want: nil,
		},
		{
			name: "Single Slot",
			args: args{
				time: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			fields: fields{
				BaseTimeSlots: []Slot{
					{
						Start: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
						End:   time.Date(2021, 1, 1, 23, 59, 59, 0, time.UTC),
						Members: []string{
							"1",
						},
					},
				},
			},
			want: &Slot{
				Start: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2021, 1, 1, 23, 59, 59, 0, time.UTC),
				Members: []string{
					"1",
				},
			},
		},
		{
			name: "Double Slot",
			args: args{
				time: time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC),
			},
			fields: fields{
				BaseTimeSlots: []Slot{
					{
						Start: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
						End:   time.Date(2021, 1, 1, 23, 59, 59, 0, time.UTC),
						Members: []string{
							"1", "2",
						},
					},
					{
						Start: time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC),
						End:   time.Date(2021, 1, 2, 23, 59, 59, 0, time.UTC),
						Members: []string{
							"2", "3",
						},
					},
				},
			},
			want: &Slot{
				Start: time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2021, 1, 2, 23, 59, 59, 0, time.UTC),
				Members: []string{
					"2", "3",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			planning := &Planning{
				BaseTimeSlots:    tt.fields.BaseTimeSlots,
				PrimaryTimeSlots: tt.fields.PrimaryTimeSlots,
			}
			if got := planning.GetActiveSlot(tt.args.time); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetActiveSlot() = %v, want %v", got, tt.want)
			}
		})
	}
}

func sameStringSlice(x, y []string) bool {
	if len(x) != len(y) {
		return false
	}
	// create a map of string -> int
	diff := make(map[string]int, len(x))
	for _, _x := range x {
		// 0 value for int is 0, so just increment a counter for the string
		diff[_x]++
	}
	for _, _y := range y {
		// If the string _y is not in diff bail out early
		if _, ok := diff[_y]; !ok {
			return false
		}
		diff[_y] -= 1
		if diff[_y] == 0 {
			delete(diff, _y)
		}
	}
	return len(diff) == 0
}
