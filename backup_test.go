package graphdb

import (
	"bytes"
	"context"
	"go-graphdb/testenv"
	"io"
	"testing"
)

func TestBackupClient_CreateRestoreLocal(t *testing.T) {
	testenv.WithEnv(t, func(url string) {
		client := New(url)

		repo, err := createRepository(t, client)
		if err != nil {
			t.Fatal(err)
		}

		backup := bytes.NewBuffer(nil)
		var target string
		var size int64
		err = client.Backups().Create(context.Background(), BackupOptions{}, func(filename string, reader io.Reader) error {
			target = filename
			size, err = io.Copy(backup, reader)
			return err
		})

		if err != nil {
			t.Fatal(err)
		}

		if target == "" {
			t.Error("expected a filename")
		}

		if size == 0 {
			t.Error("expected backup size to be > 0")
		}

		if err = client.Repositories().Delete(context.Background(), repo); err != nil {
			t.Fatal(err)
		}

		if err = client.Backups().Restore(context.Background(), RestoreOptions{}, backup); err != nil {
			t.Error("unable to restore backup")
		}
	})
}

//func TestBackupClient_CreateRestoreS3Compatible(t *testing.T) {
//	testenv.WithEnv(t, func(url string) {
//		client := New(url)
//
//		repo, err := createRepository(t, client)
//		if err != nil {
//			t.Fatal(err)
//		}
//
//		err = client.Backups().Create(context.Background(), BackupOptions{BucketUri: "s3://minio:9000/graphdb/pishki.tar?AWS_ACCESS_KEY_ID=YHZFK8M40WOBP5X3TY7R&AWS_SECRET_ACCESS_KEY=iJ412cAxX1E2qXAq2LanHNHeueFwPw1HjbeDEt2p"}, nil)
//		if err != nil {
//			t.Fatal(err)
//		}
//
//		if err = client.Repositories().Delete(context.Background(), repo); err != nil {
//			t.Fatal(err)
//		}
//
//		if err = client.Backups().Restore(context.Background(), RestoreOptions{BucketUri: "s3://minio:9000/graphdb/pishki.tar?AWS_ACCESS_KEY_ID=YHZFK8M40WOBP5X3TY7R&AWS_SECRET_ACCESS_KEY=iJ412cAxX1E2qXAq2LanHNHeueFwPw1HjbeDEt2p"}, nil); err != nil {
//			t.Error("unable to restore backup")
//		}
//	})
//}
