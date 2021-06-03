package gen

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"go/format"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/ghodss/yaml"
	"github.com/go-openapi/spec"
	"github.com/swaggo/swag"
)

var open = os.Open

// DefaultOverridesFile is the location swagger will look for type overrides.
const DefaultOverridesFile