package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func TestGetPlanning(t *testing.T) {
	type args struct {
		schedule Schedule
		date     time.Time
	}
	tests := []struct {
		name string
		args args
		want *Planning
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
			want: &Planning{
				BaseTimeSlots: []Slot{
					{
						Start: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
						End:   time.Date(2021, 1, 1, 23, 59, 59, 0, time.UTC),
						Members: []string{
							"alice", "bob",
						},
					},
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
			if got := GetPlanning(tt.args.schedule, tt.args.date); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetPlanning() = %v, want %v", got, tt.want)
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
					Id:   "1",
					Name: "alice",
				},
			},
		},
		{
			"Two Users",
			&[]User{
				{
					Id:   "1",
					Name: "alice",
				},
				{
					Id:   "2",
					Name: "bob",
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
					w.WriteHeader(http.StatusOK)
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

func TestPlanning_GetEnd(t *testing.T) {
	type fields struct {
		BaseTimeSlots    []Slot
		PrimaryTimeSlots []Slot
	}
	tests := []struct {
		name   string
		fields fields
		want   time.Time
	}{
		{
			name: "Single Slot",
			fields: fields{
				BaseTimeSlots: []Slot{
					{
						Start: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
						End:   time.Date(2021, 1, 1, 23, 59, 59, 0, time.UTC),
						Members: []string{
							"alice", "bob",
						},
					},
				},
			},
			want: time.Date(2021, 1, 1, 23, 59, 59, 0, time.UTC),
		},
		{
			name: "Double Slot",
			fields: fields{
				BaseTimeSlots: []Slot{
					{
						Start: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
						End:   time.Date(2021, 1, 1, 23, 59, 59, 0, time.UTC),
						Members: []string{
							"alice", "bob",
						},
					},
					{
						Start: time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC),
						End:   time.Date(2021, 1, 2, 23, 59, 59, 0, time.UTC),
						Members: []string{
							"alice", "bob",
						},
					},
				},
			},
			want: time.Date(2021, 1, 2, 23, 59, 59, 0, time.UTC),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			planning := &Planning{
				BaseTimeSlots:    tt.fields.BaseTimeSlots,
				PrimaryTimeSlots: tt.fields.PrimaryTimeSlots,
			}
			if got := planning.GetEnd(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetEnd() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPlanning_GetMembers(t *testing.T) {
	type fields struct {
		BaseTimeSlots    []Slot
		PrimaryTimeSlots []Slot
	}
	type args struct {
		users *[]User
	}
	users := &[]User{
		{
			Id: "1",
			Name: "Alice",
		},
		{
			Id: "2",
			Name: "Bob",
		},
		{
			Id: "3",
			Name: "Clare",
		},
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []string
	}{
		{
			name: "Single Slot",
			args: args{
				users: users,
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
			want: []string{"Alice"},
		},
		{
			name: "Double Slot",
			args: args{
				users: users,
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
			want: []string{"Alice", "Bob", "Clare"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			planning := &Planning{
				BaseTimeSlots:    tt.fields.BaseTimeSlots,
				PrimaryTimeSlots: tt.fields.PrimaryTimeSlots,
			}
			if got := planning.GetMembers(tt.args.users); !sameStringSlice(got, tt.want) {
				t.Errorf("GetMembers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPlanning_GetStart(t *testing.T) {
	type fields struct {
		BaseTimeSlots    []Slot
		PrimaryTimeSlots []Slot
	}
	tests := []struct {
		name   string
		fields fields
		want   time.Time
	}{
		{
			name: "Single Slot",
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
			want: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "Double Slot",
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
			want: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			planning := &Planning{
				BaseTimeSlots:    tt.fields.BaseTimeSlots,
				PrimaryTimeSlots: tt.fields.PrimaryTimeSlots,
			}
			if got := planning.GetStart(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetStart() = %v, want %v", got, tt.want)
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
						Start: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
						End:   time.Date(2021, 1, 1, 23, 59, 59, 0, time.UTC),
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