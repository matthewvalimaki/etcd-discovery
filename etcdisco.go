package etcdisco

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
	"text/template"
)

const (
	LOCAL_IP = "LOCAL_IP"
)

type BindingStruct struct {
	LOCAL_IP string
}

// Run a wrapped version of etcd where certain placeholder arguments
// are translated to actual values based on discovery/introspection.
func RunWrappedEtcd() {

	// get the arguments except for the target binary
	args := os.Args[1:]
	log.Printf("Original args: %v", args)

	// discover the bindings between variables and their
	// actual values (will be a map with things like local-ip -> 10.1.50.5)
	bindings, err := discoverBindings()
	log.Printf("Discovered bindings: %+v", bindings)
	if err != nil {
		log.Fatal(err)
	}

	// get a slice with post-transformation arguments
	transformedArgs, err := transformArgs(args, bindings)
	if err != nil {
		log.Fatal(err)
	}

	// invoke etcd with transformed arguments
	err = invokeEtcd(transformedArgs)
	if err != nil {
		log.Fatal(err)
	}

}

func transformArgs(args []string, bindings map[string]string) (transformed []string, err error) {

	// loop over args
	for _, arg := range args {

		// loop over bindings
		for bindingKey, bindingVal := range bindings {
			// does the argument contain bindingKey?
			if strings.Contains(arg, bindingKey) {
				// perform transformation
				transformedArg, err := transformArg(arg, bindingKey, bindingVal)
				if err != nil {
					return transformed, err
				}

				// append to output slice
				transformed = append(transformed, transformedArg)

			} else {
				// no transformation necessary, just add it to output slice
				transformed = append(transformed, arg)
			}
		}

	}

	return transformed, nil
}

// Given:
//   arg: "http://{{ local-ip }}:2379"
//   bindingKey: "local-ip"
//   bindingVal: "10.1.1.1"
// Return:
//   "http://10.1.1.1:2379"
func transformArg(arg, bindingKey, bindingVal string) (string, error) {

	tmpl, err := template.New("transform_arg").Parse(arg)
	if err != nil {
		return arg, err
	}

	params := BindingStruct{}
	switch bindingKey {
	case LOCAL_IP:
		params.LOCAL_IP = bindingVal
	}

	out := &bytes.Buffer{}

	// execute template and write to dest
	err = tmpl.Execute(out, params)
	if err != nil {
		return arg, err
	}

	return string(out.Bytes()), nil

}

func discoverBindings() (bindings map[string]string, err error) {

	bindings = map[string]string{}
	localIp, err := discoverLocalIp()
	if err != nil {
		return bindings, err
	}
	log.Printf("localIp: %v", localIp)
	bindings[LOCAL_IP] = localIp
	return bindings, nil
}

func discoverLocalIp() (localIp string, err error) {

	host, _ := os.Hostname()
	addrs, _ := net.LookupIP(host)
	for _, addr := range addrs {
		if ipv4 := addr.To4(); ipv4 != nil {
			return fmt.Sprintf("%v", ipv4), nil
		}
	}
	return "", fmt.Errorf("Could not find localip")

}

func invokeEtcd(transformedArgs []string) error {

	log.Printf("Invoking etcd with transformed args: %v", transformedArgs)
	cmd := exec.Command("etcd", transformedArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		return err
	}
	return cmd.Wait()

}
