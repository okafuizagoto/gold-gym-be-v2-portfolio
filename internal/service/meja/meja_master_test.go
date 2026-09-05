package goldgym

import (
	"context"
	"reflect"
	"testing"

	goldMejaEntity "gold-gym-be/internal/entity/meja"
	jaegerLog "gold-gym-be/pkg/log"

	"github.com/opentracing/opentracing-go"
)

func TestGenerateSequentialNames(t *testing.T) {
	cases := []struct {
		name    string
		start   string
		count   int
		want    []string
		wantErr bool
	}{
		{name: "simple", start: "A1", count: 3, want: []string{"A1", "A2", "A3"}},
		{name: "zero-pad preserved", start: "A01", count: 3, want: []string{"A01", "A02", "A03"}},
		{name: "overflow width prints wider, not truncated", start: "A98", count: 5,
			want: []string{"A98", "A99", "A100", "A101", "A102"}},
		{name: "multi-letter prefix with digits", start: "MEJA010", count: 2, want: []string{"MEJA010", "MEJA011"}},
		{name: "no trailing digit rejected", start: "VIP", count: 3, wantErr: true},
		{name: "empty start rejected", start: "   ", count: 3, wantErr: true},
		{name: "count below 1 rejected", start: "A1", count: 0, wantErr: true},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := GenerateSequentialNames(c.start, c.count)
			if c.wantErr {
				if err == nil {
					t.Fatalf("expected error, got names: %v", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("got %v, want %v", got, c.want)
			}
		})
	}
}

// fakeMejaData implements Data for InsertMejaBulk tests.
type fakeMejaData struct {
	existingNames map[string]bool
	insertedRows  []goldMejaEntity.Meja
}

func (f *fakeMejaData) InsertMeja(ctx context.Context, rows []goldMejaEntity.Meja) error {
	f.insertedRows = rows
	return nil
}

func (f *fakeMejaData) GetMejaByOutlet(ctx context.Context, goldid int, outcode string) ([]goldMejaEntity.Meja, error) {
	return nil, nil
}

func (f *fakeMejaData) GetExistingMejaNames(ctx context.Context, outcode string) (map[string]bool, error) {
	return f.existingNames, nil
}

func (f *fakeMejaData) UpdateMejaStatus(ctx context.Context, outcode string, mejaIDs []int, fromStatus, toStatus string) (int64, error) {
	return int64(len(mejaIDs)), nil
}

func (f *fakeMejaData) ReserveMeja(ctx context.Context, outcode string, mejaIDs []int) (int64, error) {
	return int64(len(mejaIDs)), nil
}

func TestInsertMejaBulk_RejectsDuplicateAgainstExisting(t *testing.T) {
	fake := &fakeMejaData{existingNames: map[string]bool{"A1": true}}
	svc := New(fake, opentracing.NoopTracer{}, jaegerLog.Factory{})

	err := svc.InsertMejaBulk(context.Background(), 1, "OUT01", []goldMejaEntity.InsertMeja{
		{MejaName: "A1", MejaCapacity: 4, MejaAreaID: 1},
	})
	if err == nil {
		t.Fatal("expected error for duplicate name against existing meja, got nil")
	}
}

func TestInsertMejaBulk_RejectsDuplicateWithinBatch(t *testing.T) {
	fake := &fakeMejaData{existingNames: map[string]bool{}}
	svc := New(fake, opentracing.NoopTracer{}, jaegerLog.Factory{})

	err := svc.InsertMejaBulk(context.Background(), 1, "OUT01", []goldMejaEntity.InsertMeja{
		{MejaName: "A1", MejaCapacity: 4, MejaAreaID: 1},
		{MejaName: "A1", MejaCapacity: 4, MejaAreaID: 1},
	})
	if err == nil {
		t.Fatal("expected error for duplicate name within same batch, got nil")
	}
}

func TestInsertMejaBulk_SuccessSetsStatusKosong(t *testing.T) {
	fake := &fakeMejaData{existingNames: map[string]bool{}}
	svc := New(fake, opentracing.NoopTracer{}, jaegerLog.Factory{})

	err := svc.InsertMejaBulk(context.Background(), 1, "OUT01", []goldMejaEntity.InsertMeja{
		{MejaName: "A1", MejaCapacity: 4, MejaAreaID: 1},
		{MejaName: "A2", MejaCapacity: 4, MejaAreaID: 1},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(fake.insertedRows) != 2 {
		t.Fatalf("expected 2 rows inserted, got %d", len(fake.insertedRows))
	}
	for _, r := range fake.insertedRows {
		if r.MejaStatus != goldMejaEntity.MejaStatusKosong {
			t.Fatalf("expected status KOSONG on insert, got %s", r.MejaStatus)
		}
	}
}
