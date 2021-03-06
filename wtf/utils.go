package wtf

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/rivo/tview"
)

const SimpleDateFormat = "Jan 2"
const SimpleTimeFormat = "15:04 MST"
const FriendlyDateTimeFormat = "Mon, Jan 2, 15:04"
const TimestampFormat = "2006-01-02T15:04:05-0700"

func CenterText(str string, width int) string {
	return fmt.Sprintf("%[1]*s", -width, fmt.Sprintf("%[1]*s", (width+len(str))/2, str))
}

func ExecuteCommand(cmd *exec.Cmd) string {
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Sprintf("%v\n", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Sprintf("%v\n", err)
	}

	var str string
	if b, err := ioutil.ReadAll(stdout); err == nil {
		str += string(b)
	}

	err = cmd.Wait()
	if err != nil {
		return fmt.Sprintf("%v\n", err)
	}

	return str
}

func Exclude(strs []string, val string) bool {
	for _, str := range strs {
		if val == str {
			return false
		}
	}
	return true
}

func FindMatch(pattern string, data string) [][]string {
	r := regexp.MustCompile(pattern)
	return r.FindAllStringSubmatch(data, -1)
}

func NameFromEmail(email string) string {
	parts := strings.Split(email, "@")
	return strings.Title(strings.Replace(parts[0], ".", " ", -1))
}

func NamesFromEmails(emails []string) []string {
	names := []string{}

	for _, email := range emails {
		names = append(names, NameFromEmail(email))
	}

	return names
}

// OpenFile opens the file defined in `path` via the operating system
func OpenFile(path string) {
	filePath, _ := ExpandHomeDir(path)
	cmd := exec.Command("open", filePath)

	ExecuteCommand(cmd)
}

// PadRow returns a padding for a row to make it the full width of the containing widget.
// Useful for ensurig row highlighting spans the full width (I suspect tcell has a better
// way to do this, but I haven't yet found it)
func PadRow(offset int, max int) string {
	padSize := max - offset
	if padSize < 0 {
		padSize = 0
	}

	return strings.Repeat(" ", padSize)
}

func ReadFileBytes(filePath string) ([]byte, error) {
	fileData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return []byte{}, err
	}

	return fileData, nil
}

func RightAlignFormat(view *tview.TextView) string {
	_, _, w, _ := view.GetInnerRect()
	return fmt.Sprintf("%%%ds", w-1)
}

func RowColor(module string, idx int) string {
	evenKey := fmt.Sprintf("wtf.mods.%s.colors.row.even", module)
	oddKey := fmt.Sprintf("wtf.mods.%s.colors.row.odd", module)

	if idx%2 == 0 {
		return Config.UString(evenKey, "white")
	} else {
		return Config.UString(oddKey, "lightblue")
	}
}

func SigilStr(len, pos int, view *tview.TextView) string {
	sigils := ""

	if len > 0 {
		sigils = strings.Repeat(Config.UString("wtf.paging.pageSigil", "*"), pos)
		sigils = sigils + Config.UString("wtf.paging.selectedSigil", "_")
		sigils = sigils + strings.Repeat(Config.UString("wtf.paging.pageSigil", "*"), len-1-pos)

		sigils = "[lightblue]" + fmt.Sprintf(RightAlignFormat(view), sigils) + "[white]"
	}

	return sigils
}

/* -------------------- Slice Conversion -------------------- */

func ToInts(slice []interface{}) []int {
	results := []int{}
	for _, val := range slice {
		results = append(results, val.(int))
	}

	return results
}

func ToStrs(slice []interface{}) []string {
	results := []string{}
	for _, val := range slice {
		results = append(results, val.(string))
	}

	return results
}

/* -------------------- Date/Time -------------------- */

// DateFormat defines the format we expect to receive dates from BambooHR in
const DateFormat = "2006-01-02"
const TimeFormat = "15:04"

func IsToday(date time.Time) bool {
	now := Now()

	return (date.Year() == now.Year()) &&
		(date.Month() == now.Month()) &&
		(date.Day() == now.Day())
}

func Now() time.Time {
	return time.Now().Local()
}

func PrettyDate(dateStr string) string {
	newTime, _ := time.Parse(DateFormat, dateStr)
	return fmt.Sprint(newTime.Format("Jan 2, 2006"))
}

func Tomorrow() time.Time {
	return Now().AddDate(0, 0, 1)
}

func UnixTime(unix int64) time.Time {
	return time.Unix(unix, 0)
}
