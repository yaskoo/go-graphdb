package graphdb

//func TestCluster_Create(t *testing.T) {
//	client := New("http://first.127.0.0.1.nip.io:7201")
//
//	_, err := client.Cluster().Config(context.Background())
//	var ce ClusterError
//	if !errors.As(err, &ce) || !ce.NotFound {
//		t.Errorf("expected no cluster error, got: %v", err)
//	}
//
//	_, err = client.Cluster().Create(context.Background(), ClusterConfig{
//		Nodes: []string{"graphdb.example.com"},
//	})
//
//	if !errors.As(err, &ce) || !ce.Messages.IsList() || len(ce.Messages.List()) != 2 {
//		t.Errorf("expected two error messages, got %v", ce.Messages.List())
//	}
//
//	err = client.Cluster().Truncate(context.Background())
//
//	if !errors.As(err, &ce) || len(ce.Messages.List()) != 1 {
//		t.Errorf("expected truncate log to fail because there is no cluster")
//	}
//
//	_, err = client.Cluster().Create(context.Background(), ClusterConfig{
//		Nodes: []string{
//			"first.127.0.0.1.nip.io:7301",
//			"second.127.0.0.1.nip.io:8000",
//		},
//	})
//
//	if !errors.As(err, &ce) || len(ce.Messages.Map()) != 1 {
//		t.Errorf("expected one error message, got %v", ce.Messages.List())
//	}
//
//	_, err = client.Cluster().AddNodes(context.Background(), []string{"second.127.0.0.1.nip.io:7302"})
//	if !errors.As(err, &ce) || ce.Messages.IsList() == false || len(ce.Messages.List()) != 1 {
//		t.Error("expected add nodes to fail, because there is no cluster")
//	}
//
//	_, err = client.Cluster().DeleteNodes(context.Background(), []string{"second.127.0.0.1.nip.io:7302"})
//	if !errors.As(err, &ce) || ce.Messages.IsList() == false || len(ce.Messages.List()) != 1 {
//		t.Error("expected delete nodes to fail, because there is no cluster")
//	}
//
//	_, err = client.Cluster().ReplaceNodes(context.Background(), []string{"fourth.127.0.0.1.nip.io:7304"}, []string{"second.127.0.0.1.nip.io:7302"})
//	if !errors.As(err, &ce) || ce.Messages.IsList() == false || len(ce.Messages.List()) != 1 {
//		t.Error("expected replace nodes to fail, because there is no cluster")
//	}
//
//	nodes := []string{
//		"first.127.0.0.1.nip.io:7301",
//		"second.127.0.0.1.nip.io:7302",
//		"third.127.0.0.1.nip.io:7303",
//	}
//	re, err := client.Cluster().Create(context.Background(), ClusterConfig{
//		Nodes: nodes,
//	})
//
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	if len(re) != len(nodes) {
//		t.Fatalf("expected %d success messages, got %d", len(nodes), len(re))
//	}
//
//	status, err := client.Cluster().NodeStatus(context.Background())
//	if err != nil {
//		t.Errorf("error getting node status: %v", err)
//	}
//
//	if status.Endpoint == "" {
//		t.Errorf("expected cluster to be enabled")
//	}
//
//	statuses, err := client.Cluster().Status(context.Background())
//	if err != nil {
//		t.Errorf("error getting cluster status: %v", err)
//	}
//
//	if len(statuses) != len(nodes) {
//		t.Errorf("expected %d node statuses, got %d", len(nodes), len(statuses))
//	}
//
//	config, err := client.Cluster().Config(context.Background())
//	if err != nil {
//		t.Errorf("error getting cluster config: %v", err)
//	}
//
//	if config.ElectionRangeTimeout == 0 {
//		t.Error("expected non-zero cluster config")
//	}
//
//	config, err = client.Cluster().Update(context.Background(), ClusterProperties{TransactionLogMaximumSizeGB: 5.0})
//	if err != nil {
//		t.Errorf("error updating cluster config: %v", err)
//	}
//
//	if config.TransactionLogMaximumSizeGB != 5.0 {
//		t.Errorf("expected transaction log maximum size to be 5, got %f", config.TransactionLogMaximumSizeGB)
//	}
//
//	re, err = client.Cluster().DeleteNodes(context.Background(), []string{"second.127.0.0.1.nip.io:7302"})
//	if err != nil {
//		t.Errorf("error deleting nodes: %v", err)
//	}
//
//	re, err = client.Cluster().AddNodes(context.Background(), []string{"second.127.0.0.1.nip.io:7302"})
//	if err != nil {
//		t.Errorf("error adding nodes: %v", err)
//	}
//
//	if len(re) != 1 {
//		t.Errorf("expected %d success messages, got %d", 1, len(re))
//	}
//
//	re, err = client.Cluster().ReplaceNodes(context.Background(), []string{"fourth.127.0.0.1.nip.io:7304"}, []string{"second.127.0.0.1.nip.io:7302"})
//	if err != nil {
//		t.Errorf("error replacing nodes: %v", err)
//	}
//
//	if len(re) != 2 {
//		t.Errorf("expected %d success messages, got %d", 2, len(re))
//	}
//
//	err = client.Cluster().Truncate(context.Background())
//	if err != nil {
//		t.Errorf("error truncating cluster log: %v", err)
//	}
//
//	err = client.Cluster().AddTag(context.Background(), "east")
//	if err != nil {
//		t.Errorf("error adding tag: %v", err)
//	}
//
//	err = client.Cluster().RemoveTag(context.Background(), "east")
//	if err != nil {
//		t.Errorf("error removing tag: %v", err)
//	}
//
//	err = client.Cluster().Delete(context.Background())
//	if err != nil {
//		t.Errorf("error deleting cluster: %v", err)
//	}
//}
