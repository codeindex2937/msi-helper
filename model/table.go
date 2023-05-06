package msi

type ActionText struct {
	Action      string     `msi:"Action,key"`
	Description NullString `msi:"Description"`
	Template    NullString `msi:"Template"`
}

type AdminExecuteSequence struct {
	Action    string     `msi:"Action,key"`
	Condition NullString `msi:"Condition"`
	Sequence  NullInt16  `msi:"Sequence"`
}

type AdminUISequence struct {
	Action    string     `msi:"Action,key"`
	Condition NullString `msi:"Condition"`
	Sequence  NullInt16  `msi:"Sequence"`
}

type AdvtExecuteSequence struct {
	Action    string     `msi:"Action,key"`
	Condition NullString `msi:"Condition"`
	Sequence  NullInt16  `msi:"Sequence"`
}

type AppSearch struct {
	Property  string `msi:"Property,key"`
	Signature string `msi:"Signature_,key"`
}

type BBControl struct {
	Billboard  string     `msi:"Billboard_,key"`
	BBControl  string     `msi:"BBControl,key"`
	Type       string     `msi:"Type"`
	X          int16      `msi:"X"`
	Y          int16      `msi:"Y"`
	Width      int16      `msi:"Width"`
	Height     int16      `msi:"Height"`
	Attributes NullInt32  `msi:"Attributes"`
	Text       NullString `msi:"Text"`
}

type Billboard struct {
	Billboard string     `msi:"Billboard,key"`
	Feature   string     `msi:"Feature_"`
	Action    NullString `msi:"Action"`
	Ordering  NullInt16  `msi:"Ordering"`
}

type Binary struct {
	Name string `msi:"Name,key"`
	Data Stream `msi:"Data"`
}

type CheckBox struct {
	Property string     `msi:"Property,key"`
	Value    NullString `msi:"Value"`
}

type ComboBox struct {
	Property string     `msi:"Property,key"`
	Order    int16      `msi:"Order,key"`
	Value    string     `msi:"Value"`
	Text     NullString `msi:"Text"`
}

type Component struct {
	Component   string     `msi:"Component,key"`
	ComponentId NullString `msi:"ComponentId"`
	Directory   string     `msi:"Directory_"`
	Attributes  int16      `msi:"Attributes"`
	Condition   NullString `msi:"Condition"`
	KeyPath     NullString `msi:"KeyPath"`
	Feature     string
	Files       map[string]FileCredential
}

type Control struct {
	Dialog      string     `msi:"Dialog_,key"`
	Control     string     `msi:"Control,key"`
	Type        string     `msi:"Type"`
	X           int16      `msi:"X"`
	Y           int16      `msi:"Y"`
	Width       int16      `msi:"Width"`
	Height      int16      `msi:"Height"`
	Attributes  NullInt32  `msi:"Attributes"`
	Property    NullString `msi:"Property"`
	Text        NullString `msi:"Text"`
	ControlNext NullString `msi:"Control_Next"`
	Help        NullString `msi:"Help"`
}

type ControlCondition struct {
	Dialog    string `msi:"Dialog_,key"`
	Control   string `msi:"Control_,key"`
	Action    string `msi:"Action,key"`
	Condition string `msi:"Condition,key"`
}

type ControlEvent struct {
	Dialog    string     `msi:"Dialog_,key"`
	Control   string     `msi:"Control_,key"`
	Event     string     `msi:"Event,key"`
	Argument  string     `msi:"Argument,key"`
	Condition NullString `msi:"Condition,key"`
	Ordering  NullInt16  `msi:"Ordering"`
}

type CreateFolder struct {
	Directory string `msi:"Directory_,key"`
	Component string `msi:"Component_,key"`
}

type CustomAction struct {
	Action       string     `msi:"Action,key"`
	Type         int16      `msi:"Type"`
	Source       NullString `msi:"Source"`
	Target       NullString `msi:"Target"`
	ExtendedType NullInt32  `msi:"ExtendedType"`
}

type Dialog struct {
	Dialog         string     `msi:"Dialog,key"`
	HCentering     int16      `msi:"HCentering"`
	VCentering     int16      `msi:"VCentering"`
	Width          int16      `msi:"Width"`
	Height         int16      `msi:"Height"`
	Attributes     NullInt32  `msi:"Attributes"`
	Title          NullString `msi:"Title"`
	ControlFirst   string     `msi:"Control_First"`
	ControlDefault NullString `msi:"Control_Default"`
	ControlCancel  NullString `msi:"Control_Cancel"`
}

type Directory struct {
	Directory       string     `msi:"Directory,key"`
	DirectoryParent NullString `msi:"Directory_Parent"`
	DefaultDir      string     `msi:"DefaultDir"`
	Dirs            []*Directory
}

type DrLocator struct {
	Signature string     `msi:"Signature_,key"`
	Parent    NullString `msi:"Parent,key"`
	Path      NullString `msi:"Path,key"`
	Depth     NullInt16  `msi:"Depth"`
}

type Error struct {
	Error   int16      `msi:"Error,key"`
	Message NullString `msi:"Message"`
}

type EventMapping struct {
	Dialog    string `msi:"Dialog_,key"`
	Control   string `msi:"Control_,key"`
	Event     string `msi:"Event,key"`
	Attribute string `msi:"Attribute"`
}

type Feature struct {
	Feature       string     `msi:"Feature,key"`
	FeatureParent NullString `msi:"Feature_Parent"`
	Title         NullString `msi:"Title"`
	Description   NullString `msi:"Description"`
	Display       NullInt16  `msi:"Display"`
	Level         int16      `msi:"Level"`
	Directory     NullString `msi:"Directory_"`
	Attributes    int16      `msi:"Attributes"`
}

type FeatureComponents struct {
	Feature   string `msi:"Feature_,key"`
	Component string `msi:"Component_,key"`
}

type File struct {
	File       string     `msi:"File,key"`
	Component  string     `msi:"Component_"`
	FileName   string     `msi:"FileName"`
	FileSize   int32      `msi:"FileSize"`
	Version    NullString `msi:"Version"`
	Language   NullString `msi:"Language"`
	Attributes NullInt16  `msi:"Attributes"`
	Sequence   int16      `msi:"Sequence"`
}

type Icon struct {
	Name string `msi:"Name,key"`
	Data Stream `msi:"Data"`
}

type InstallExecuteSequence struct {
	Action    string     `msi:"Action,key"`
	Condition NullString `msi:"Condition"`
	Sequence  NullInt16  `msi:"Sequence"`
}

type InstallUISequence struct {
	Action    string     `msi:"Action,key"`
	Condition NullString `msi:"Condition"`
	Sequence  NullInt16  `msi:"Sequence"`
}

type LaunchCondition struct {
	Condition   string `msi:"Condition,key"`
	Description string `msi:"Description"`
}

type ListBox struct {
	Property string     `msi:"Property,key"`
	Order    int16      `msi:"Order,key"`
	Value    string     `msi:"Value"`
	Text     NullString `msi:"Text"`
}

type ListView struct {
	Property string     `msi:"Property,key"`
	Order    int16      `msi:"Order,key"`
	Value    string     `msi:"Value"`
	Text     NullString `msi:"Text"`
	Binary   NullString `msi:"Binary_"`
}

type Media struct {
	DiskId       int16      `msi:"DiskId,key"`
	LastSequence int16      `msi:"LastSequence"`
	DiskPrompt   NullString `msi:"DiskPrompt"`
	Cabinet      NullString `msi:"Cabinet"`
	VolumeLabel  NullString `msi:"VolumeLabel"`
	Source       NullString `msi:"Source"`
}

type Property struct {
	Property string `msi:"Property,key"`
	Value    string `msi:"Value"`
}

type RadioButton struct {
	Property string     `msi:"Property,key"`
	Order    int16      `msi:"Order,key"`
	Value    string     `msi:"Value"`
	X        int16      `msi:"X"`
	Y        int16      `msi:"Y"`
	Width    int16      `msi:"Width"`
	Height   int16      `msi:"Height"`
	Text     NullString `msi:"Text"`
	Help     NullString `msi:"Help"`
}

type RegLocator struct {
	Signature string     `msi:"Signature_,key"`
	Root      int16      `msi:"Root"`
	Key       string     `msi:"Key"`
	Name      NullString `msi:"Name"`
	Type      NullInt16  `msi:"Type"`
}

type Registry struct {
	Registry  string     `msi:"Registry,key"`
	Root      int16      `msi:"Root"`
	Key       string     `msi:"Key"`
	Name      NullString `msi:"Name"`
	Value     NullString `msi:"Value"`
	Component string     `msi:"Component_"`
}

type RemoveFile struct {
	FileKey     string     `msi:"FileKey,key"`
	Component   string     `msi:"Component_"`
	FileName    NullString `msi:"FileName"`
	DirProperty string     `msi:"DirProperty"`
	InstallMode int16      `msi:"InstallMode"`
}

type RemoveRegistry struct {
	RemoveRegistry string     `msi:"RemoveRegistry,key"`
	Root           int16      `msi:"Root"`
	Key            string     `msi:"Key"`
	Name           NullString `msi:"Name"`
	Component      string     `msi:"Component_"`
}

type ServiceControl struct {
	ServiceControl string     `msi:"ServiceControl,key"`
	Name           string     `msi:"Name"`
	Event          int16      `msi:"Event"`
	Arguments      NullString `msi:"Arguments"`
	Wait           NullInt16  `msi:"Wait"`
	Component      string     `msi:"Component_"`
}

type ServiceInstall struct {
	ServiceInstall string     `msi:"ServiceInstall,key"`
	Name           string     `msi:"Name"`
	DisplayName    NullString `msi:"DisplayName"`
	ServiceType    int32      `msi:"ServiceType"`
	StartType      int32      `msi:"StartType"`
	ErrorControl   int32      `msi:"ErrorControl"`
	LoadOrderGroup NullString `msi:"LoadOrderGroup"`
	Dependencies   NullString `msi:"Dependencies"`
	StartName      NullString `msi:"StartName"`
	Password       NullString `msi:"Password"`
	Arguments      NullString `msi:"Arguments"`
	Component      string     `msi:"Component_"`
	Description    NullString `msi:"Description"`
}

type Shortcut struct {
	Shortcut               string     `msi:"Shortcut,key"`
	Directory              string     `msi:"Directory_"`
	Name                   string     `msi:"Name"`
	Component              string     `msi:"Component_"`
	Target                 string     `msi:"Target"`
	Arguments              NullString `msi:"Arguments"`
	Description            NullString `msi:"Description"`
	Hotkey                 NullInt16  `msi:"Hotkey"`
	Icon                   NullString `msi:"Icon_"`
	IconIndex              NullInt16  `msi:"IconIndex"`
	ShowCmd                NullInt16  `msi:"ShowCmd"`
	WkDir                  NullString `msi:"WkDir"`
	DisplayResourceDLL     NullString `msi:"DisplayResourceDLL"`
	DisplayResourceId      NullInt32  `msi:"DisplayResourceId"`
	DescriptionResourceDLL NullString `msi:"DescriptionResourceDLL"`
	DescriptionResourceId  NullInt32  `msi:"DescriptionResourceId"`
}

type Signature struct {
	Signature  string     `msi:"Signature,key"`
	FileName   string     `msi:"FileName"`
	MinVersion NullString `msi:"MinVersion"`
	MaxVersion NullString `msi:"MaxVersion"`
	MinSize    NullInt32  `msi:"MinSize"`
	MaxSize    NullInt32  `msi:"MaxSize"`
	MinDate    NullInt32  `msi:"MinDate"`
	MaxDate    NullInt32  `msi:"MaxDate"`
	Languages  NullString `msi:"Languages"`
}

type TextStyle struct {
	TextStyle string    `msi:"TextStyle,key"`
	FaceName  string    `msi:"FaceName"`
	Size      int16     `msi:"Size"`
	Color     NullInt32 `msi:"Color"`
	StyleBits NullInt16 `msi:"StyleBits"`
}

type UIText struct {
	Key  string     `msi:"Key,key"`
	Text NullString `msi:"Text"`
}

type Upgrade struct {
	UpgradeCode    string     `msi:"UpgradeCode,key"`
	VersionMin     NullString `msi:"VersionMin,key"`
	VersionMax     NullString `msi:"VersionMax,key"`
	Language       NullString `msi:"Language,key"`
	Attributes     int32      `msi:"Attributes,key"`
	Remove         NullString `msi:"Remove"`
	ActionProperty string     `msi:"ActionProperty"`
}

type Validation struct {
	Table       string     `msi:"Table,key"`
	Column      string     `msi:"Column,key"`
	Nullable    string     `msi:"Nullable"`
	MinValue    NullInt32  `msi:"MinValue"`
	MaxValue    NullInt32  `msi:"MaxValue"`
	KeyTable    NullString `msi:"KeyTable"`
	KeyColumn   NullInt16  `msi:"KeyColumn"`
	Category    NullString `msi:"Category"`
	Set         NullString `msi:"Set"`
	Description NullString `msi:"Description"`
}

type Streams struct {
	Name string `msi:"Name,key"`
	Data Stream `msi:"Data"`
}
