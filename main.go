package main

import (
	"fmt"
	"github.com/andlabs/ui"
	_ "github.com/andlabs/ui/winmanifest"
	"math"
	"regexp"
	"strconv"
	"strings"
)

var (
	unidade = []string{"", "um", "dois", "três", "quatro", "cinco",
		"seis", "sete", "oito", "nove", "dez", "onze",
		"doze", "treze", "quatorze", "quinze", "dezesseis",
		"dezessete", "dezoito", "dezenove"}
	centena = []string{"", "cento", "duzentos", "trezentos",
		"quatrocentos", "quinhentos", "seiscentos",
		"setecentos", "oitocentos", "novecentos"}
	dezena = []string{"", "", "vinte", "trinta", "quarenta", "cinquenta",
		"sessenta", "setenta", "oitenta", "noventa"}
	qualificaS = []string{"", "mil", "milhão", "bilhão", "trilhão"}
	qualificaP = []string{"", "mil", "milhões", "bilhões", "trilhões"}

	w   *ui.Window
)

func main() {
	ui.Main(setupUI)
}

func setupUI() {
	w = ui.NewWindow(fmt.Sprintf("Número/Rais por Extenso"), 640, 480, false)
	w.OnClosing(func(*ui.Window) bool {
		ui.Quit()
		return true
	})
	ui.OnShouldQuit(func() bool {
		w.Destroy()
		return true
	})

	vbox := ui.NewVerticalBox()


	entryForm := ui.NewForm()
	entryForm.SetPadded(true)
	vbox.Append(entryForm, true)

	multilineEntry := ui.NewMultilineEntry()
	multilineEntry.SetReadOnly(true)

	entry := ui.NewEntry()
	entry.OnChanged(func(entry *ui.Entry) {
		if entry.Text() == "" {
			multilineEntry.SetText("")
			return
		}

		compile := regexp.MustCompile(`(\d)|(\.)|(,)`)
		subMatch := compile.FindAllString(entry.Text(), -1)

		text := strings.Join(subMatch, "")
		text = strings.Replace(text, "..", ".", -1)
		text = strings.Replace(text, ",,", ",", -1)
		entry.SetText(text)

		split := strings.Split(entry.Text(), ",")

		compileReais := regexp.
			MustCompile(`^((\d{1,3}\.\d{3}\.\d{3}\.\d{3})|(\d{1,3}\.\d{3}\.\d{3})|(\d{1,3}\.\d{3})|(\d+))$`)
		if !compileReais.MatchString(split[0]) {
			multilineEntry.SetText("Erro: valor não permitido.")
			return
		}

		reais := split[0]
		reais = strings.Replace(reais, ".", "", -1)

		cents := "00"
		compileCents := regexp.MustCompile(`(\d,$)|(,\d{1,2}$)`)
		if len(split) > 1 {
			if !compileCents.MatchString(entry.Text()) || len(split) > 2 {
				multilineEntry.SetText("Erro: valor não permitido.")
				return
			}
			cents = split[1]
		}

		aux, _ := strconv.Atoi(fmt.Sprintf("%s%s", reais, cents))
		number := float64(aux) / 100
		multilineEntry.SetText(ValueInFull(number))
	})


	entryForm.Append("Número", entry, false)	
	entryForm.Append("por Extenso", multilineEntry, true)

	w.SetChild(vbox)
	w.SetMargined(true)

	w.Show()
}

func ValueInFull(value float64) string {
	var (
		vlrP   string
		s      string
		saux   string
		n      int
		unid   int
		dez    int
		cent   int
		tam    int
		i      = 0
		umReal bool
		tem    bool
	)

	if value == 0.0 {
		return "zero"
	}

	inteiro := int(value)
	resto := (value - float64(inteiro)) * 100.0

	if inteiro > 999999999 {
		return "Erro: valor superior a 999 trilhões."
	}

	vlrS := strconv.Itoa(inteiro)

	centavos := strconv.Itoa(int(math.Round(resto)))

	for vlrS != "0" {
		tam = len(vlrS)
		// retira do valor a 1a. parte, 2a. parte, por exemplo, para 123456789:
		// 1a. parte = 789 (centena)
		// 2a. parte = 456 (mil)
		// 3a. parte = 123 (milhões)
		if tam > 3 {
			vlrP = string([]rune(vlrS)[tam-3:])
			vlrS = string([]rune(vlrS)[:tam-3])
		} else { // última parte do valor
			vlrP = vlrS
			vlrS = "0"
		}

		if vlrP != "000" {
			saux = ""
			if vlrP == "100" {
				saux = "cem"
			} else {
				n, _ = strconv.Atoi(vlrP) // para n = 371, tem-se:
				cent = n / 100            // cent = 3 (centena trezentos)
				dez = (n % 100) / 10      // dez  = 7 (dezena setenta)
				unid = (n % 100) % 10     // unid = 1 (unidade um)
				if cent != 0 {
					saux = centena[cent]
				}
				if (n % 100) <= 19 {
					if len(saux) != 0 {
						saux = saux + " e " + unidade[n%100]
					} else {
						saux = unidade[n%100]
					}
				} else {
					if len(saux) != 0 {
						saux = saux + " e " + dezena[dez]
					} else {
						saux = dezena[dez]
					}
					if unid != 0 {
						if len(saux) != 0 {
							saux = saux + " e " + unidade[unid]
						} else {
							saux = unidade[unid]
						}
					}
				}
			}
			if vlrP == "1" || vlrP == "001" {
				if i == 0 { // 1a. parte do valor (um real)
					umReal = true
				} else {
					saux = saux + " " + qualificaS[i]
				}
			} else if i != 0 {
				saux = saux + " " + qualificaP[i]
			}
			if len(s) != 0 {
				s = saux + ", " + s
			} else {
				s = saux
			}
		}
		if (i == 0 || i == 1) && len(s) != 0 {
			tem = true // tem centena ou mil no valor
		}
		i = i + 1 // próximo qualificador: 1- mil, 2- milhão, 3- bilhão, ...
	}

	if len(s) != 0 {
		if umReal {
			s = s + " real"
		} else if tem {
			s = s + " reais"
		} else {
			s = s + " de reais"
		}
	}

	// definindo o extenso dos centavos do valor
	if centavos != "0" { // valor com centavos
		if len(s) != 0 { // se não é valor somente com centavos
			s = s + " e "
		}
		if centavos == "1" {
			s = s + "um centavo"
		} else {
			n, _ = strconv.Atoi(centavos)
			if n <= 19 {
				s = s + unidade[n]
			} else { // para n = 37, tem-se:
				unid = n % 10 // unid = 37 % 10 = 7 (unidade sete)
				dez = n / 10  // dez  = 37 / 10 = 3 (dezena trinta)
				s = s + dezena[dez]
				if unid != 0 {
					s = s + " e " + unidade[unid]
				}
			}
			s = s + " centavos"
		}
	}
	return s
}
