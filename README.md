## TODO

Post task to /task/name
Struct for result

Commands :
PS
```go
package main

import (
	"fmt"
	ps "github.com/mitchellh/go-ps"
)

func main() {
	processList, err := ps.Processes()
	if err != nil {
		return
	}
	// map ages
	for x := range processList {
		var process ps.Process
		process = processList[x]
		fmt.Printf("%5d\t%5d\t%s\n", process.PPid(), process.Pid(), process.Executable())
	}
}

```
WMI-AV
```go
package main

import (
	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
)

func main() {
	// init COM, oh yeah
	err := ole.CoInitialize(0)
	if err != nil {
		return
	}
	defer ole.CoUninitialize()

	unknown, _ := oleutil.CreateObject("WbemScripting.SWbemLocator")
	defer unknown.Release()

	wmi, _ := unknown.QueryInterface(ole.IID_IDispatch)
	defer wmi.Release()

	// service is a SWbemServices
	serviceRaw, _ := oleutil.CallMethod(wmi, "ConnectServer", nil, "root/SecurityCenter2")
	service := serviceRaw.ToIDispatch()
	defer service.Release()

	// result is a SWBemObjectSet
	resultRaw, _ := oleutil.CallMethod(service, "ExecQuery", "SELECT * FROM AntiVirusProduct")
	result := resultRaw.ToIDispatch()
	defer result.Release()

	countVar, _ := oleutil.GetProperty(result, "Count")
	count := int(countVar.Val)

	for i := 0; i < count; i++ {
		itemRaw, _ := oleutil.CallMethod(result, "ItemIndex", i)
		item := itemRaw.ToIDispatch()
		asString, _ := oleutil.GetProperty(item, "displayName")
		println(asString.ToString())
	}
}
```
https://fourcore.io/blogs/manipulating-windows-tokens-with-golang