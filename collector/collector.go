package collector

import (
	"sync"
	"time"
)

type CollectDevice struct {
	Time        time.Time
	SpaceStats  map[string]SpaceStatFS `json:"spaceStats"`
	NvidiaStats []SmiNvidia            `json:"nvidiaStats"`
}

type SpaceStatFS struct {
	ReservedBlocks uint64 `json:"reservedBlocks"`
	BlockSize      uint64 `json:"blockSize"`
	Total          uint64 `json:"total"`
	Free           uint64 `json:"free"`
}

type SmartDisk struct {
	Model       string `json:"modelFamily"`
	Capacity    string `json:"capacity"`
	Health      string `json:"health"`
	Used        string `json:"used"`
	Working     int64  `json:"working"`
	Temperature string `json:"temperature"`
	Full        string `json:"full"`
	Error       string `json:"error"`
}

type RaidMD struct {
	Disks   []string  `json:"disks"`
	Proc    RaidProc  `json:"proc"`
	ProcOut string    `json:"procOut"`
	Adm     RaidMDAdm `json:"adm"`
	AdmOut  string    `json:"admOut"`
}

type RaidProc struct {
	State    string  `json:"state"`
	Progress float64 `json:"progress"`
	Left     float64 `json:"left"`
	Speed    int     `json:"speed"`
	ParseErr string  `json:"parseErr,omitempty"`
}

type RaidMDAdm struct {
	Name      string `json:"name"`
	State     string `json:"state"`
	Level     string `json:"level"`
	Capacity  string `json:"capacity"`
	CreatedAt string `json:"createdAt"`
	Active    string `json:"active"`
	Working   string `json:"working"`
	Failed    string `json:"failed"`
	Spare     string `json:"spare"`
}

type DisksInfo struct {
	Version int                   `json:"version"`
	Time    time.Time             `json:"time"`
	Smarts  map[string]*SmartDisk `json:"smarts"`
	Raids   map[string]*RaidMD    `json:"raids"`
}

type SmiNvidia struct {
	Name       string `json:"name"`
	Driver     string `json:"driver"`
	PciDev     string `json:"pciDev"`
	PciBus     string `json:"pciBus"`
	PciLink    string `json:"pciLink"`
	PciTx      string `json:"pciTx"`
	PciRx      string `json:"pciRx"`
	Fan        string `json:"fan"`
	MemUse     string `json:"memUse"`
	MemFree    string `json:"memFree"`
	UtilGpu    string `json:"utilGpu"`
	UtilMem    string `json:"utilMem"`
	UtilEnc    string `json:"utilEnc"`
	UtilDec    string `json:"utilDec"`
	TmpGpu     string `json:"tmpGpu"`
	TmpMem     string `json:"tmpMem"`
	Power      string `json:"power"`
	PowerLimit string `json:"powerLimit"`
	ClockGra   string `json:"clockSh"`
	ClockSm    string `json:"clockSm"`
	ClockMem   string `json:"clockMem"`
	ClockVideo string `json:"clockVideo"`
}

type Collector struct {
	mu   sync.Mutex
	data CollectDevice

	ChanDisksInfo chan *DisksInfo
}

func New() *Collector {
	c := &Collector{
		mu:            sync.Mutex{},
		data:          CollectDevice{},
		ChanDisksInfo: make(chan *DisksInfo, 1),
	}

	go c.collectSpace()
	go c.collectDisks()
	go c.collectGpu()

	return c
}

func (c *Collector) Get() CollectDevice {
	return c.data
}
