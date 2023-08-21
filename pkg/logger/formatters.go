package logger

import (
	"context"
	"path/filepath"
	"runtime"
	"time"
)

const (
	// Escape цветовые последовательности для раскраски лога (не везде работет всё)
	// EscStart + bold; + colorSym; + colorBack + EscEnd @see Wikipedia.org

	EscStart    = "\x1b["
	EscColorEnd = "m"

	EscReset       = "0"    // 	Reset / Normal 	выключение всех атрибутов
	EscBold        = "1"    // 	Жирный или увеличить яркость
	EscDark        = "2"    // 	Блёклый (уменьшить яркость) 	Не везде поддерживается
	EscCoursive    = "3"    // 	Курсив: вкл. 	Не везде поддерживается. Иногда обрабатывается как инверсия.
	EscUnderLine   = "4"    // 	Подчёркнутый: один раз
	EscBlinkSlow   = "5"    //	Мигание: Медленно 	менее 150 раз в минуту
	EscBlinkFast   = "6"    // 	Мигание: Часто 	MS-DOS ANSI.SYS; 150+ в минуту; не везде поддерживается
	EscInverse     = "7"    // 	Отображение: Негатив 	инвертирует или обращает; меняет цвета фона и текста
	EscHide        = "8"    // 	Скрытый 	Не везде поддерживается.
	EscDeprecate   = "9"    // 	Зачёркнутый 	Символы разборчивы, но помечены как удалённые. Не везде поддерживается.
	EscNormalFont  = "10"   // 	Основной (по умолчанию) шрифт
	EscFacture     = "20"   // 	Фрактура 	вряд ли поддерживается
	EscNoBold      = "21"   // 	Жирный: выкл. или Подчёркивание: Двойное 	жирный не всегда поддерживается; двойное подчёркивание вряд ли поддерживается.
	EscNormalLight = "22"   // 	Обычный цвет или яркость 	Ни жирный ни блеклый
	EscNoCursive   = "23"   // 	Не курсивный, не фрактура
	EscNoUnderLine = "24"   // 	Подчёркивание: Нет 	Подчёркивание ни одиночное ни двойное
	EscNoBlink     = "25"   // 	Мигание: выкл.
	EscNoIverse    = "27"   // 	Отображение: обычное 	не негатив
	EscNoHide      = "28"   // 	Отображающийся 	выключить скрытие
	EscNoDeprecate = "29"   // 	Не зачёркнутый
	Esc256Sym      = "38;5" //	Зарезервировано для дополнительных цветов 5;n где n {\displaystyle n} n индекс цвета (0..255)
	EscRGBSym      = "38;2" // или 2;r;g;b где r , g , b {\displaystyle r,g,b} {\displaystyle r,g,b} — красный, зелёный и синий каналы цвета (в пределах 255)
	EscDefSym      = "39"   // 	Цвет текста по умолчанию (на переднем плане) 	зависит от реализации (в соответствии со стандартом)
	Esc256Back     = "48;5" // 	Зарезервировано для установки расширенного цвета фона 5;n где n {\displaystyle n} n индекс цвета (0..255)
	EscRGBBack     = "48;2" // или 2;r;g;b где r , g , b {\displaystyle r,g,b} {\displaystyle r,g,b} — красный, зелёный и синий каналы цвета (в пределах 255)
	EscDefBack     = "49"   // 	Цвет фона по умолчанию 	зависит от реализации (в соответствии со стандартом)
	EscBorder1     = "51"   // 	Обрамлённый
	EscBorder2     = "52"   // 	Окружённый
	EscUpperLine   = "53"   // 	Надчёркнутый
	EscNoBorder    = "54"   // 	Не обрамлённый и не окружённый
	EscNoUpperLine = "55"   // 	Не надчёркнутый

	// 11–19: альтернативный шрифт с 0 по 9
	// 30–37: цвет текста (в скобках: +EscBold) черный(серый), красный(ярко..), зеленый(ярко..), болотный(желтый),
	//               синий(голубой), маджента(светлее), циан(светлее), серый(белый)
	// 40–47: аналогично цвет фона

	// Интегральная сборка типовых раскрасок:

	EscBlackCurrent     = "\x1b[10;30m" // черным по текущему (черный == hide!)
	EscRedCurrent       = "\x1b[10;31m" // темнокрасный ..
	EscGreenCurrent     = "\x1b[10;32m" // темнозеленый по текущему
	EscYellowCurrent    = "\x1b[10;33m" // коричневый(?) по текущему
	EscBlueCurrent      = "\x1b[10;34m" // синим по текущему
	EscMagentaCurrent   = "\x1b[10;35m" // фиолетовым по текущему
	EscCyanCurrent      = "\x1b[10;36m" // сианом по текущему
	EscLightGrayCurrent = "\x1b[10;37m" // светосерым по текущему
	EscWhiteCurrent     = "\x1b[1;37m"  // белым(ярким, жирным) по текущему
)

// FormatTime -- добавляет в заданный буфер время в текущем формате лога
func FormatTime(buf *[]byte, now time.Time, flags int) {
	if flags&LogUTC != 0 {
		//		Now = Now.UTC()
	}
	if flags&LogDate != 0 {
		y, m, d := now.Date()
		*buf = append(*buf, ' ')
		itoaBuf(buf, y, 4)
		*buf = append(*buf, '-')
		itoaBuf(buf, int(m), 2)
		*buf = append(*buf, '-')
		itoaBuf(buf, d, 2)
	}
	if flags&(LogTime|LogMicroSeconds) != 0 {
		h, m, s := now.Clock()
		*buf = append(*buf, ' ')
		itoaBuf(buf, h, 2)
		*buf = append(*buf, ':')
		itoaBuf(buf, m, 2)
		*buf = append(*buf, ':')
		itoaBuf(buf, s, 2)

		if flags&LogMicroSeconds != 0 {
			*buf = append(*buf, '.')
			itoaBuf(buf, now.Nanosecond()/1e3, 3)
		}
	}

}

// FormatFileLine -- добавляет в буфер информацию о файле и номере строки
func FormatFileLine(buf *[]byte, depth int, isShort bool) {
	var (
		ok   bool
		line int
		file string
	)
	if _, file, line, ok = runtime.Caller(depth); !ok {
		file = "???"
		line = 0
	}
	if isShort {
		file = filepath.Base(file)
	}
	*buf = append(*buf, file...)
	*buf = append(*buf, '#')
	itoaBuf(buf, line, 4) // max 9999 line number!
}

// FormatFuncLine -- добавляет в буфер название функции(метода) и номер строки файла
func FormatFuncLine(buf *[]byte, depth int) {
	var frame runtime.Frame
	var ok bool

	if frame, ok = GetCaller(depth); !ok {
		frame.Function = "???"
		frame.Line = 0
	}
	*buf = append(*buf, frame.Function...)
	*buf = append(*buf, '#')
	itoaBuf(buf, frame.Line, 4) // max line = 9999 !!!
}

// FormatTrace -- достает из контекста значение по заданному иденту и добавляет его в лог
func FormatTrace(buf *[]byte, ctx context.Context, traceId any) {
	var traceVal string
	var ok bool

	traceTag := ctx.Value(traceId)
	if traceTag != nil {
		if traceVal, ok = traceTag.(string); !ok {
			traceVal = "???"
		}
	}

	*buf = append(*buf, traceVal...)
}

// FormatSlice -- добавление в строку лога заданного списка текстовок
func FormatSlice(buf *[]byte, messages []string) {
	for _, message := range messages {
		*buf = append(*buf, message...)
		*buf = append(*buf, ',')
		*buf = append(*buf, ' ')
	}
}

// GlDefColors -- дефолтная таблица цветов для вывода сообщения
var GlDefColors = map[string]string{
	LogPanicPrefix: EscStart + EscBlinkFast + EscColorEnd + EscStart + EscBold + ";" + EscRedCurrent + EscColorEnd,
	LogFatalPrefix: EscStart + EscBlinkSlow + EscColorEnd + EscStart + EscBold + ";" + EscRedCurrent + EscColorEnd,
	LogErrorPrefix: EscStart + EscBold + ";" + EscRedCurrent + EscColorEnd,
	LogWarnPrefix:  EscStart + EscNormalFont + ";" + EscMagentaCurrent + EscColorEnd,
	LogInfoPrefix:  EscStart + EscNormalFont + ";" + EscGreenCurrent + EscColorEnd,
	LogDebugPrefix: EscStart + EscCoursive + ";" + EscBlueCurrent + EscColorEnd,
}

// FormatColored -- раскраска строки заданным цветом с последующей отменой
func FormatColored(buf *[]byte, color, message string) {
	*buf = append(*buf, color...)
	*buf = append(*buf, message...)
	*buf = append(*buf, EscStart+EscReset+EscColorEnd...)
}
