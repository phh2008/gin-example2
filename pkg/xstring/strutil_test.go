package xstring

import (
	"log"
	"testing"
)

func TestSnakeCase(t *testing.T) {
	//EquipmentMonitoringProject
	//equipment_monitoring_project
	var words = "EquipmentMonitoringProject.CpuGhz"
	des := SnakeCase(words)
	log.Println(des)
}
