package nilpointer

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"golang.org/x/tools/go/analysis"
	"sort"
	"strings"
)

var Analyzer = &analysis.Analyzer{
	Name: "npecheck",
	Doc:  Doc,
	Run:  Run,
}

const Doc = "check potential nil pointer reference"

// CheckPointerPosition 记录检查过空指针的代码位置
type CheckPointerPosition struct {
	Line      int
	Colum     int
	IsChecked bool
	Type      int // DefaultPtrType, SlicePtrType, ParentPtrCurNonType
}

const (
	DefaultPtrType      int = 0 // ptr
	SlicePtrType        int = 1 // []ptr
	ParentPtrCurNonType int = 2 //  A.B.C , A 是指针，当前B 非指针
)

// IsPointer 判断一个类型是否是指针类型
func IsPointer(typ types.Type) bool {
	_, ok := typ.(*types.Pointer)
	return ok
}

func IsSliceIncludePointerElem(typ types.Type) bool {
	var sliceElem, isSlice = typ.(*types.Slice)
	if isSlice {
		elem := sliceElem.Elem()
		if _, ok := elem.(*types.Pointer); ok {
			return true
		}
	}
	return false
}

func GetNodeType(ident *ast.Ident, typesInfo *types.Info) (int, bool) {
	var (
		nodeType           = NodeTypeDefaultSinglePtr
		isReturnSingleFunc = false
	)
	obj := typesInfo.ObjectOf(ident) // 需要加以判断函数的单个返回值是不是ptr

	// 可以尝试看看是不是函数, 并且函数要有一个返回值
	sign, ok := obj.Type().(*types.Signature)
	if ok && sign != nil && sign.Results() != nil && sign.Results().Len() == 1 { // 函数、方法
		isReturnSingleFunc = true
		retType := sign.Results().At(0).Type()
		if !IsPointer(retType) {
			nodeType = NodeTypeNonSinglePtr
		}
	} else { // 变量
		if !IsPointer(obj.Type()) {
			nodeType = NodeTypeNonSinglePtr
		}
	}

	return nodeType, isReturnSingleFunc
}

// IsPointerArray 这个和上面的IsSliceIncludePointerElem 入参不一样,目的一样
func IsPointerArray(ident *ast.Ident, info *types.Info) bool {
	obj := info.ObjectOf(ident)
	if obj == nil {
		// ident 不是一个有效的标识符
		return false
	}

	typ := obj.Type()
	arr, ok := typ.(*types.Slice)
	if !ok {
		// ident 不是一个数组类型
		return false
	}

	elem := arr.Elem()
	_, ok = elem.(*types.Pointer)
	if ok {
		// 数组元素是一个指针类型
		return true
	}

	//if ptr.Elem().String() == targetType {
	//	// 这是一个指向 targetType 类型的指针数组
	//	return true
	//}

	return false
}

func GetIdentPosition(p *token.Position, ident *ast.Ident, fset *token.FileSet) {
	*p = fset.Position(ident.Pos())
}

func WalkSelector(expr *ast.SelectorExpr, fset *token.FileSet, walkFunc func(ast.Node)) {
	walkFunc(expr)

	ident, ok := expr.X.(*ast.Ident)
	if ok {
		walkFunc(ident)
		return
	}

	if se, ok := expr.X.(*ast.SelectorExpr); ok {
		WalkSelector(se, fset, walkFunc)
	}
}

func (f *FuncDelChecker) isRootComeFromOutside(varName string) bool {
	if varName == "" {
		return false
	}

	varNameNodes := strings.Split(varName, ".")
	if len(varNameNodes) >= 1 {
		if _, ok := f.needCheckPointerPositionMap[varNameNodes[0]]; ok {
			return true
		}
	}

	return false
}

func (f *FuncDelChecker) isExistFuncRetPtr(nodeList []*SelectNode, isNeedRemoveLeaf bool) (bool, int) {
	// 存在函数返回指针，比如 o.GetUserInfo().GetMembership()
	if len(nodeList) == 0 {
		return false, -1
	}

	if isNeedRemoveLeaf {
		nodeList = nodeList[0 : len(nodeList)-1] // 去除叶子节点
	}

	for index, n := range nodeList {
		if n.Type == NodeTypeDefaultSinglePtr && n.IsReturnSingleFunc == true {
			return true, index
		}
	}

	return false, -2
}

func (f *FuncDelChecker) recordIfBinaryNilValidation(binaryExpr *ast.BinaryExpr, fset *token.FileSet, lintErrors *[]*LintError) {
	if binaryExpr == nil {
		return
	}
	y := binaryExpr.Y
	if ident, ok := y.(*ast.Ident); ok {
		if ident.Name != "nil" { // 后面可以看看这里有无其他项目种常用的alias, 比如RedisNil
			return
		}

		if binaryExpr.Op != token.NEQ && binaryExpr.Op != token.EQL { // 只要 != nil , == nil 判断
			return
		}
		x := binaryExpr.X
		if expr, ok := x.(*ast.Ident); ok { // 例如 if a != nil 走这个分支
			var pos token.Position
			GetIdentPosition(&pos, expr, fset) // 记录判段 nil 的位置
			// nodeType, isReturnSingleFunc := GetNodeType(expr, f.pass.TypesInfo)
			if f.isRootComeFromOutside(expr.Name) /*|| (isReturnSingleFunc == true && nodeType == NodeTypeDefaultSinglePtr) */ {
				checkedPos := &CheckPointerPosition{
					Line:      pos.Line,
					Colum:     pos.Column,
					IsChecked: true,
				}

				index := f.findFirstSuitablePosIndexFromStart(expr.Name, pos)
				if index >= 0 {
					if f.needCheckPointerPositionMap[expr.Name][index].IsChecked != true { // 状态一样就不用去改line colum 了
						f.needCheckPointerPositionMap[expr.Name][index] = checkedPos

						if originName, ok := f.originFieldMap[expr.Name]; ok {
							f.needCheckPointerPositionMap[originName] = f.needCheckPointerPositionMap[expr.Name]
						}
					}
				}
			}
		}
	}

	binaryExpr = f.recordBinaryExpr(binaryExpr, fset, lintErrors)
}

func (f *FuncDelChecker) recordBinaryExpr(
	expr *ast.BinaryExpr,
	fset *token.FileSet,
	lintErrorList *[]*LintError) *ast.BinaryExpr {
	left := expr.X
	right := expr.Y
	// 如果左右两侧有子 BinaryExpr，则递归处理子 BinaryExpr
	switch leftExpr := left.(type) {
	case *ast.BinaryExpr:
		x := leftExpr.X
		switch expr := x.(type) {
		case *ast.SelectorExpr: // 例如 if a.b != nil // if a.b.c == nil 走这个分支
			_ = f.travelSelectorNameAndFuncWithRecord(expr, fset, lintErrorList)
		case *ast.Ident: // 例如 if a != nil 走这个分支
			f.recordIdentCheckedPosition(expr, fset)
		}
		f.recordBinaryExpr(leftExpr, fset, lintErrorList)

	case *ast.SelectorExpr: // 例如 if a.b != nil // if a.b.c == nil 走这个分支
		_ = f.travelSelectorNameAndFuncWithRecord(leftExpr, fset, lintErrorList)

	case *ast.CallExpr:
		if selectExpr, ok := leftExpr.Fun.(*ast.SelectorExpr); ok {
			_ = f.travelSelectorNameAndFuncWithRecord(selectExpr, fset, lintErrorList)
		}
	}

	switch rightExpr := right.(type) {
	case *ast.BinaryExpr:
		x := rightExpr.X
		switch expr := x.(type) {
		case *ast.SelectorExpr: // 例如 if a.b != nil // if a.b.c == nil 走这个分支
			_ = f.travelSelectorNameAndFuncWithRecord(expr, fset, lintErrorList)

		case *ast.Ident: // 例如 if a != nil 走这个分支
			f.recordIdentCheckedPosition(expr, fset)
		}
		f.recordBinaryExpr(rightExpr, fset, lintErrorList)

	case *ast.SelectorExpr: // 例如 if a.b != nil // if a.b.c == nil 走这个分支
		_ = f.travelSelectorNameAndFuncWithRecord(rightExpr, fset, lintErrorList)

	case *ast.CallExpr:
		if selectExpr, ok := rightExpr.Fun.(*ast.SelectorExpr); ok {
			_ = f.travelSelectorNameAndFuncWithRecord(selectExpr, fset, lintErrorList)
		}
	}

	// 新建一个 BinaryExpr
	newExpr := &ast.BinaryExpr{
		Op: expr.Op,
		X:  left,
		Y:  right,
	}

	return newExpr
}

func (f *FuncDelChecker) isExistRecord(name string) ([]*CheckPointerPosition, bool) {
	recordList, ok := f.needCheckPointerPositionMap[name]
	return recordList, ok
}

// 初始化的时候，遇到同名name ,是append, 这里需要找到做if nil 判断的时候，记录该记在哪个index
// 写检测位置用
func (f *FuncDelChecker) findFirstSuitablePosIndexFromStart(name string, pos token.Position) int {
	needCheckPointerPositionList, ok := f.needCheckPointerPositionMap[name]
	if !ok || len(needCheckPointerPositionList) == 0 {
		return -1 // name不存在,  新起记录还是不记录，上层去判断
	}

	if len(needCheckPointerPositionList) == 1 { // 无重复name
		return 0
	}

	for i, v := range needCheckPointerPositionList { // 有重复var name
		if pos.Line > v.Line || (pos.Line == v.Line && pos.Column > v.Colum) { // 大于其一
			// 小于下一个 ， 或者没有下一个
			nextIndex := i + 1
			if nextIndex == len(needCheckPointerPositionList) {
				return i
			}

			nextV := needCheckPointerPositionList[nextIndex]
			if pos.Line < nextV.Line || (pos.Line == v.Line && pos.Column < nextV.Colum) {
				return i
			}
		}
	}

	return -2 // 上层可以选择append
}

// 读检测位置用
func (f *FuncDelChecker) findFirstSuitablePosIndexFromEnd(name string, pos token.Position) int {
	needCheckPointerPositionList, ok := f.needCheckPointerPositionMap[name]
	if !ok {
		return -1 // 不需要检测，name不存在
	}

	index := len(needCheckPointerPositionList) - 1
	for index >= 0 {
		if needCheckPointerPositionList[index].Line < pos.Line ||
			(needCheckPointerPositionList[index].Line == pos.Line &&
				needCheckPointerPositionList[index].Colum < pos.Column) {
			return index
		}

		index--
	}

	return -2 // Name存在，但不存在对应的Pos
}

func (f *FuncDelChecker) recordIdentCheckedPosition(expr *ast.Ident, fset *token.FileSet) {
	var pos token.Position
	GetIdentPosition(&pos, expr, fset)      // 记录判段 nil 的位置
	if f.isRootComeFromOutside(expr.Name) { // 只用于限制外部变量
		checkedPosition := &CheckPointerPosition{
			Line:      pos.Line,
			Colum:     pos.Column,
			IsChecked: true,
		}

		// 找到合适的位置，去记录做过检测的位置
		index := f.findFirstSuitablePosIndexFromStart(expr.Name, pos)
		if index >= 0 {
			if f.needCheckPointerPositionMap[expr.Name][index].IsChecked != true {
				f.needCheckPointerPositionMap[expr.Name][index] = checkedPosition

				// 若有alias , 把原始的也记录一下检测
				if originName, ok := f.originFieldMap[expr.Name]; ok {
					f.needCheckPointerPositionMap[originName] = f.needCheckPointerPositionMap[expr.Name]
				}
			}
		}
	}
}

func (f *FuncDelChecker) preRecordNilPointerFromOutside(fnDel *ast.FuncDecl, lintErrors *[]*LintError) {
	if fnDel == nil {
		return
	}

	var (
		typeInfo = f.pass.TypesInfo
		fset     = f.pass.Fset
	)

	// 获取函数入参为指针的变量
	if fnDel.Type != nil && fnDel.Type.Params != nil {
		for _, field := range fnDel.Type.Params.List {
			for _, name := range field.Names {
				// 函数参数 ，记录参数指针
				typ := typeInfo.Types[field.Type].Type
				var pos token.Position
				GetIdentPosition(&pos, name, fset)
				if IsPointer(typ) {
					f.needCheckPointerPositionMap[name.Name] = []*CheckPointerPosition{
						{
							Line:      pos.Line,
							Colum:     pos.Column,
							IsChecked: false, // 入参一进来的时候还没有检查过
						},
					}
				} else if IsSliceIncludePointerElem(typ) {
					f.needCheckPointerPositionMap[name.Name] = []*CheckPointerPosition{
						{
							Line:      pos.Line,
							Colum:     pos.Column,
							IsChecked: true, // slice 这一层不用ptr校验，下面的指针元素需要校验
							Type:      SlicePtrType,
						},
					}
				}
			}
		}
	}

	// 获取函数外部赋值的指针变量
	for _, stmt := range fnDel.Body.List {
		f.recordStmtNilValidation(stmt, fset, lintErrors, typeInfo)
	}
}

func (f *FuncDelChecker) recordStmtNilValidation(stmt ast.Stmt, fset *token.FileSet, lintErrors *[]*LintError, typeInfo *types.Info) {
	switch s := stmt.(type) {
	case *ast.IfStmt: // 记录做过if nil 判断的位置
		f.recordIfStmtNilValidation(s, fset, lintErrors, typeInfo)

	case *ast.AssignStmt: // 赋值语句, 识别外部函数赋值
		f.recordAssignmentStmtNilValidation(s, typeInfo, fset)

	case *ast.RangeStmt:
		f.recordRangeStmtNilValidation(s, typeInfo, fset, lintErrors)

	case *ast.SwitchStmt:
		f.recordSwitchStmtNilValidation(s, typeInfo, fset, lintErrors)
	}
}

func (f *FuncDelChecker) recordIfStmtNilValidation(s *ast.IfStmt, fset *token.FileSet, lintErrors *[]*LintError, typeInfo *types.Info) {
	cond := s.Cond
	switch cond := cond.(type) {
	case *ast.BinaryExpr: // 最后再考虑 if a := anotherScope(); a != nil
		f.recordIfBinaryNilValidation(cond, fset, lintErrors)
	}

	for _, stmt := range s.Body.List {
		switch expr := stmt.(type) {
		case *ast.IfStmt, *ast.AssignStmt, *ast.RangeStmt:
			f.recordStmtNilValidation(expr, fset, lintErrors, typeInfo)
		}
	}
}

func (f *FuncDelChecker) recordSwitchStmtNilValidation(s *ast.SwitchStmt, typeInfo *types.Info, fset *token.FileSet, lintErrors *[]*LintError) {
	if s == nil {
		return
	}

	if s.Body == nil {
		return
	}

	for _, c := range s.Body.List {
		if caseClause, ok := c.(*ast.CaseClause); ok {
			for _, b := range caseClause.Body {
				switch expr := b.(type) {
				case *ast.AssignStmt, *ast.IfStmt, *ast.SwitchStmt, *ast.RangeStmt:
					f.recordStmtNilValidation(expr, fset, lintErrors, typeInfo)
				}
			}
		}
	}
}

// 也会检测 for k,v := range xx.GetYYY() 右边的表达式空指针
func (f *FuncDelChecker) recordRangeStmtNilValidation(s *ast.RangeStmt, typeInfo *types.Info, fset *token.FileSet, lintErrors *[]*LintError) {
	var (
		x               = s.X
		rangeParentName string
	)

	switch expr := x.(type) {
	case *ast.Ident:
		rangeParentName = expr.Name
		if _, ok := f.needCheckPointerPositionMap[rangeParentName]; ok {
			value := s.Value
			switch value.(type) {
			case *ast.Ident:
				f.recordRangeValue(s, fset)
			}
		}

	case *ast.SelectorExpr: // 记录 for k,v := range resp.list 中的 v
		// 找到叶子节点， 如果发现叶子节点是包含指针的数组，记录
		innerX := expr.X
		switch innerX := innerX.(type) {
		case *ast.Ident:
			rangeParentName = innerX.Name
			if _, ok := f.needCheckPointerPositionMap[rangeParentName]; ok {
				if IsPointerArray(expr.Sel, f.pass.TypesInfo) {
					lintErrorList := f.detectSelectorReferenceWithFunc(expr, fset) // 比如resp.list , resp 自己就没检测过
					if len(lintErrorList) > 0 {
						*lintErrors = append(*lintErrors, lintErrorList...)
					}

					f.recordRangeValue(s, fset)
				}
			}
		}

	// userInfoList := GetUserInfoList()
	// userInfoList := xx.GetUserInfoList()
	// userInfoList := xx.yy.GetUserInfoList() 这些在表达式语句、赋值语句需要考虑
	// for k, v := range xx.yy.GetUserInfoList() 需要把v 记录起来
	case *ast.CallExpr: // userInfoList := xx.GetUserInfoList()
		if exprFun, ok := expr.Fun.(*ast.SelectorExpr); ok {
			lintErrorList := f.detectSelectorReferenceWithFunc(exprFun, fset) // 比如resp.GetXXX() , resp 自己就没检测过
			if len(lintErrorList) > 0 {
				*lintErrors = append(*lintErrors, lintErrorList...)
			}
		}

		sign := GetFuncSignature(expr, typeInfo)
		if sign.Results().Len() != 1 { //  for _ , v := range list ， list 作为一个ident
			return
		}

		retType := sign.Results().At(0).Type()
		if IsSliceIncludePointerElem(retType) {
			f.recordRangeValue(s, fset)
		}
	}

	body := s.Body
	if body == nil {
		return
	}

	for _, b := range body.List {
		switch expr := b.(type) {
		case *ast.IfStmt, *ast.AssignStmt, *ast.RangeStmt, *ast.SwitchStmt:
			f.recordStmtNilValidation(expr, fset, lintErrors, typeInfo)
		}
	}
}

func (f *FuncDelChecker) recordRangeValue(s *ast.RangeStmt, fset *token.FileSet) {
	value := s.Value
	if v, ok := value.(*ast.Ident); ok {
		var pos token.Position
		GetIdentPosition(&pos, v, fset) // 记录 list 中的指针元素
		needCheckPos := &CheckPointerPosition{
			Line:      pos.Line,
			Colum:     pos.Column,
			IsChecked: false,
		}
		recordList, ok := f.isExistRecord(v.Name) //需要检测，但没有检测过
		if ok {
			f.needCheckPointerPositionMap[v.Name] = append(recordList, needCheckPos)
		} else {
			f.needCheckPointerPositionMap[v.Name] = []*CheckPointerPosition{needCheckPos}
		}
	}
}

func (f *FuncDelChecker) recordAssignmentStmtNilValidation(s *ast.AssignStmt, typeInfo *types.Info, fset *token.FileSet) {
	// 遍历变量列表
	var respVarNameList = make([]string, 0)
	for _, l := range s.Lhs {
		if i, ok := l.(*ast.Ident); ok { // 下划线变量名，要注意 ，看看跟函数返回值个数要匹配
			respVarNameList = append(respVarNameList, i.Name)
		}
	}

	// 检测右边的Rhs 里是否有来自外部的变量X，以及是否有调用函数，函数返回值类型是否是指针，如果有则记录这个指针变量
	for _, expr := range s.Rhs {
		switch EX := expr.(type) {
		case *ast.SelectorExpr: // 例如 b := a.B ， c := a.b.c
			obj := typeInfo.ObjectOf(EX.Sel) // 叶子节点
			if IsPointer(obj.Type()) {       // 最后的Sel 是指针
				fieldNameList := TravelSelectorName(EX, fset)
				if len(fieldNameList) == 0 {
					continue
				}
				// 如果最前面的字段是函数入参传过来，或者外部函数赋值过来
				if _, ok := f.needCheckPointerPositionMap[fieldNameList[0]]; ok { // x.Name 是第一个字段名
					key := strings.Join(fieldNameList, ".")
					if len(respVarNameList) == 1 { // c := a.b 这种操作，左边只可能一个值
						f.originFieldMap[respVarNameList[0]] = key
						var pos token.Position
						GetIdentPosition(&pos, EX.Sel, fset)
						needCheckedPos := &CheckPointerPosition{
							Line:      pos.Line,
							Colum:     pos.Column,
							IsChecked: false,
						}

						recordList, ok := f.isExistRecord(respVarNameList[0])
						if ok {
							f.needCheckPointerPositionMap[respVarNameList[0]] = append(recordList, needCheckedPos)
						} else {
							f.needCheckPointerPositionMap[respVarNameList[0]] = []*CheckPointerPosition{needCheckedPos}
						}
						f.needCheckPointerPositionMap[key] = f.needCheckPointerPositionMap[respVarNameList[0]]
					}
				}
			}

		case *ast.CallExpr: // 右侧可能多个函数，但其实取一个就可以了
			// 获取 CallExpr 中的函数对象和函数签名 // 如果获取不到类型，就不用判断，大部分HTTP 请求，返回参数也是指针
			sign := GetFuncSignature(EX, typeInfo)
			if sign == nil {
				continue
			}

			// 判断函数返回类型是否为指针类型
			respVarNameListLength := len(respVarNameList)
			if respVarNameListLength > 0 && sign != nil && sign.Results() != nil && respVarNameListLength == sign.Results().Len() {
				for index, varName := range respVarNameList {
					retType := sign.Results().At(index).Type()
					pos := fset.Position(EX.Rparen)
					var (
						needCheckedPos *CheckPointerPosition
						isNeedRecord   bool
					)
					if IsSliceIncludePointerElem(retType) {
						needCheckedPos = &CheckPointerPosition{
							Line:      pos.Line,
							Colum:     pos.Column,
							IsChecked: true, // slice 层默认checked
						}
						isNeedRecord = true
					}

					if IsPointer(retType) {
						// 函数返回类型为指针类型
						needCheckedPos = &CheckPointerPosition{
							Line:      pos.Line,
							Colum:     pos.Column,
							IsChecked: false,
						}
						isNeedRecord = true
					}

					if isNeedRecord {
						recordList, ok := f.isExistRecord(varName)
						if ok {
							f.needCheckPointerPositionMap[varName] = append(recordList, needCheckedPos)
						} else {
							f.needCheckPointerPositionMap[varName] = []*CheckPointerPosition{needCheckedPos}
						}
					}
				}
			}
		}
	}
}

// 定义错误类型
type LintError struct {
	Message string
	File    string
	Line    int
	Colum   int
}

const NPEMessageTipInfo = "potential nil pointer reference"

// 实现 Error 方法
func (err *LintError) Error() string {
	return fmt.Sprintf("%s: %s:%d:%d", err.Message, err.File, err.Line, err.Colum)
}

// 去除变量叶子节点, 比如a输出a,  a.b 输出a,  a.b.c 输出a.b
func RemoveVarLeafNode(varName string, originFieldMap map[string]string) string {
	if varName == "" {
		return ""
	}

	parts := strings.Split(varName, ".") // 将字符串切割为多个部分
	var substr string
	//if len(parts) == 1 {
	//	substr = parts[0] // 只有一个节点，就不用remove 了
	//}

	if len(parts) > 1 {
		substr = strings.Join(parts[:len(parts)-1], ".")
	}

	// 若有alias ， 还原
	if originName, ok := originFieldMap[substr]; ok {
		substr = originName
	}

	return substr
}

type SelectNode struct {
	Name               string
	Type               int // NodeTypeDefaultSinglePtr, NodeTypeNonSinglePtr
	IsReturnSingleFunc bool
	CurIdent           *ast.Ident
}

const (
	NodeTypeDefaultSinglePtr int = 0
	NodeTypeNonSinglePtr     int = 1 // 非默认的单个指针, 非叶子节点，代表非指针, 叶子节点任意类型
)

func buildFieldNameFromNodes(nodes []*SelectNode) string {
	fieldNameList := make([]string, 0)
	for _, n := range nodes {
		fieldNameList = append(fieldNameList, n.Name)
	}

	return strings.Join(fieldNameList, ".")
}

func walkSelectorWithFunc(expr *ast.SelectorExpr, fset *token.FileSet, walkFunc func(ast.Node)) {
	walkFunc(expr)

	switch ex := expr.X.(type) {
	case *ast.Ident:
		walkFunc(ex)
		return

	case *ast.SelectorExpr:
		walkSelectorWithFunc(ex, fset, walkFunc)

	case *ast.CallExpr:
		if se, ok := ex.Fun.(*ast.SelectorExpr); ok {
			walkSelectorWithFunc(se, fset, walkFunc)
		}

		if se, ok := ex.Fun.(*ast.Ident); ok {
			walkFunc(se)
			return
		}
	}
}

func travelSelectorFullName(expr *ast.SelectorExpr, fset *token.FileSet, typeInfo *types.Info) []string {
	var nodeNameList []string
	walkSelectorWithFunc(expr, fset, func(node ast.Node) {
		switch exprInner := node.(type) {
		case *ast.SelectorExpr:
			if exprInner != nil && exprInner.Sel != nil {
				nodeNameList = append(nodeNameList, exprInner.Sel.Name)
			}

		case *ast.CallExpr:
			switch exFunc := exprInner.Fun.(type) {
			case *ast.SelectorExpr:
				if exFunc != nil && exFunc.Sel != nil {
					nodeNameList = append(nodeNameList, exFunc.Sel.Name)
				}

			case *ast.Ident:
				nodeNameList = append(nodeNameList, exFunc.Name)
			}

		case *ast.Ident: // 最后一个，会结束递归
			if exprInner != nil {
				nodeNameList = append(nodeNameList, exprInner.Name)
				sort.SliceStable(nodeNameList, func(i, j int) bool { // 反转数组
					return i > j
				})
			}
		}
	})

	return nodeNameList
}

// travel selector [包含检测链式函数], 重点记录检测位置, 并上报检测过程中的NPE ---reference  那个改一下，以及， 不是仅仅判断 rootFromOutside 了
func (f *FuncDelChecker) travelSelectorNameAndFuncWithRecord(
	expr *ast.SelectorExpr,
	fset *token.FileSet,
	lintErrorList *[]*LintError,
) []*SelectNode {
	var nodeList = make([]*SelectNode, 0)
	walkSelectorWithFunc(expr, fset, func(node ast.Node) {
		switch exprInner := node.(type) {
		case *ast.SelectorExpr:
			if exprInner != nil && exprInner.Sel != nil {
				nodeType, isReturnSingleFunc := GetNodeType(exprInner.Sel, f.pass.TypesInfo) // 可以再记录一个是否是函数
				nodeList = append(nodeList, &SelectNode{
					Name:               exprInner.Sel.Name,
					CurIdent:           exprInner.Sel,
					Type:               nodeType,
					IsReturnSingleFunc: isReturnSingleFunc,
				})
			}

		case *ast.CallExpr:
			switch exFunc := exprInner.Fun.(type) {
			case *ast.SelectorExpr: // 方法调用 o.GetUserInfo()
				if exFunc != nil && exFunc.Sel != nil {
					nodeType, isReturnSingleFunc := GetNodeType(exFunc.Sel, f.pass.TypesInfo)
					nodeList = append(nodeList, &SelectNode{
						Name:               exFunc.Sel.Name,
						CurIdent:           exFunc.Sel,
						Type:               nodeType,
						IsReturnSingleFunc: isReturnSingleFunc})
				}

			case *ast.Ident: // 函数调用 o.GetUserInfo().GetMembership() 的叶子节点
				nodeType, isReturnSingleFunc := GetNodeType(exFunc, f.pass.TypesInfo)
				nodeList = append(nodeList, &SelectNode{Name: exFunc.Name,
					CurIdent:           exFunc,
					Type:               nodeType,
					IsReturnSingleFunc: isReturnSingleFunc,
				})
			}

		case *ast.Ident: // 最后一个，会结束递归
			if exprInner != nil {
				nodeType, isReturnSingleFunc := GetNodeType(exprInner, f.pass.TypesInfo)
				nodeList = append(nodeList, &SelectNode{
					Name:               exprInner.Name,
					CurIdent:           exprInner,
					Type:               nodeType,
					IsReturnSingleFunc: isReturnSingleFunc})
				sort.SliceStable(nodeList, func(i, j int) bool { // 反转数组
					return i > j
				})
			}

			// 记录判段 nil 的位置
			var pos token.Position
			GetIdentPosition(&pos, exprInner, fset)
			fieldName := buildFieldNameFromNodes(nodeList)
			//_, funcIndex := f.isExistFuncRetPtr(nodeList, false)
			isRootFromOutside := f.isRootComeFromOutside(fieldName)
			if isRootFromOutside /*|| isExistFunc*/ { // root 来自外部变量，或者存在函数并且返回值是指针
				// 记录前得检测一下前面的情况 ， 比如这次是 if a.b != nil , 需要先看看 a 检测了没有
				lintErrors := f.hasSequenceDetectNode(nodeList, pos, exprInner, isRootFromOutside) // 如果是 o.GetUserInfo().GetMembership(),o是Receiver 不用检测
				if len(lintErrors) != 0 {                                                          // 前面的都没有检测， 上报error
					*lintErrorList = append(*lintErrorList, lintErrors...)
				} else {
					checkedPosition := &CheckPointerPosition{
						Line:      pos.Line,
						Colum:     pos.Column,
						IsChecked: true,
						Type:      DefaultPtrType,
					}

					leafNodeType := nodeList[len(nodeList)-1].Type
					if leafNodeType == NodeTypeNonSinglePtr {
						checkedPosition = &CheckPointerPosition{
							Line:      0,
							Colum:     0,
							IsChecked: true,
							Type:      ParentPtrCurNonType,
						}
					}

					// 找到合适的位置，去记录做过检测的位置, 数组也默认检测过了
					index := f.findFirstSuitablePosIndexFromStart(fieldName, pos)
					if index >= 0 {
						if f.needCheckPointerPositionMap[fieldName][index].IsChecked != true {
							f.needCheckPointerPositionMap[fieldName][index] = checkedPosition
						}
					} else {
						f.needCheckPointerPositionMap[fieldName] = []*CheckPointerPosition{checkedPosition}
					}

					// 若有alias , 把原始的也记录一下检测
					if originName, ok := f.originFieldMap[fieldName]; ok {
						f.needCheckPointerPositionMap[originName] = f.needCheckPointerPositionMap[fieldName]
					}
				}
			}
		}
	})

	return nodeList
}

// travel selector, 重点记录检测位置, 并上报检测过程中的NPE
//func (f *FuncDelChecker) travelSelectorNameWithRecord(
//	expr *ast.SelectorExpr,
//	fset *token.FileSet,
//	lintErrorList *[]*LintError,
//) []*SelectNode {
//	var nodeNameList = make([]*SelectNode, 0)
//	WalkSelector(expr, fset, func(node ast.Node) {
//		switch exprInner := node.(type) {
//		case *ast.SelectorExpr:
//			if exprInner != nil && exprInner.Sel != nil {
//				nodeType := GetNodeType(exprInner.Sel, f.pass.TypesInfo)
//				nodeNameList = append(nodeNameList, &SelectNode{Name: exprInner.Sel.Name, Type: nodeType})
//			}
//
//		case *ast.Ident: // 最后一个，会结束递归
//			if exprInner != nil {
//				nodeType := GetNodeType(exprInner, f.pass.TypesInfo)
//				nodeNameList = append(nodeNameList, &SelectNode{Name: exprInner.Name, Type: nodeType})
//				sort.SliceStable(nodeNameList, func(i, j int) bool { // 反转数组
//					return i > j
//				})
//			}
//
//			// 记录判段 nil 的位置
//			var pos token.Position
//			GetIdentPosition(&pos, exprInner, fset)
//			//fieldName := strings.Join(nodeNameList, ".")
//			fieldName := buildFieldNameFromNodes(nodeNameList)
//			if f.isRootComeFromOutside(fieldName) { // 只用于限制外部变量
//				// 记录前得检测一下前面的情况 ， 比如这次是 if a.b != nil , 需要先看看 a 检测了没有
//				//lintErrors := f.hasSequenceDetectNodeName(fieldName, pos, exprInner)
//				lintErrors := f.hasSequenceDetectNode(nodeNameList, pos, exprInner)
//				if len(lintErrors) != 0 { // 前面的都没有检测， 上报error
//					*lintErrorList = append(*lintErrorList, lintErrors...)
//				} else {
//					checkedPosition := &CheckPointerPosition{
//						Line:      pos.Line,
//						Colum:     pos.Column,
//						IsChecked: true,
//						Type:      DefaultPtrType,
//					}
//
//					nodeType := nodeNameList[len(nodeNameList)-1].Type
//					if nodeType == NodeTypeNonSinglePtr {
//						checkedPosition = &CheckPointerPosition{
//							Line:      0,
//							Colum:     0,
//							IsChecked: true,
//							Type:      ParentPtrCurNonType,
//						}
//					}
//
//					// 找到合适的位置，去记录做过检测的位置, 数组也默认检测过了
//					index := f.findFirstSuitablePosIndexFromStart(fieldName, pos)
//					if index >= 0 && f.needCheckPointerPositionMap[fieldName][index].IsChecked != true {
//						f.needCheckPointerPositionMap[fieldName][index] = checkedPosition
//					} else {
//						f.needCheckPointerPositionMap[fieldName] = []*CheckPointerPosition{checkedPosition}
//					}
//
//					// 若有alias , 把原始的也记录一下检测
//					if originName, ok := f.originFieldMap[fieldName]; ok {
//						f.needCheckPointerPositionMap[originName] = f.needCheckPointerPositionMap[fieldName]
//					}
//				}
//			}
//		}
//	})
//
//	return nodeNameList
//}

// 递归获取完整的 selector name
func TravelSelectorName(expr *ast.SelectorExpr, fset *token.FileSet) []string {
	var nodeNameList []string
	WalkSelector(expr, fset, func(node ast.Node) {
		switch exprInner := node.(type) {
		case *ast.SelectorExpr:
			if exprInner != nil && exprInner.Sel != nil {
				nodeNameList = append(nodeNameList, exprInner.Sel.Name)
			}

		case *ast.Ident: // 最后一个，会结束递归
			if exprInner != nil {
				nodeNameList = append(nodeNameList, exprInner.Name)
				sort.SliceStable(nodeNameList, func(i, j int) bool { // 反转数组
					return i > j
				})
			}
		}
	})

	return nodeNameList
}

// 是否潜在的空指针引用, 如果是, 先汇总信息，最后再一起报错panic
func (f *FuncDelChecker) getPotentialNilPointerReference(
	varName string,
	filePath string,
	pos *token.Position,
	expr *ast.Ident,
	isComeFromOutSide bool,
) *LintError {
	if pos == nil || len(f.needCheckPointerPositionMap) == 0 {
		return nil
	}

	refName := RemoveVarLeafNode(varName, f.originFieldMap)
	if refName == "" {
		return nil
	}

	//needCheckInfo, ok := needCheckPointerPositionMap[refName]
	index := f.findFirstSuitablePosIndexFromEnd(refName, *pos)

	if isComeFromOutSide && index < 0 {
		f.pass.Reportf(expr.Pos(), NPEMessageTipInfo)
		return &LintError{
			Message: NPEMessageTipInfo,
			File:    filePath,
			Line:    pos.Line,
			Colum:   pos.Column,
		}
	}

	if index >= 0 {
		if lintError := f.buildLintError(f.needCheckPointerPositionMap[refName][index], filePath, pos); lintError != nil {
			f.pass.Reportf(expr.Pos(), NPEMessageTipInfo)
			return lintError
		}
	}

	return nil
}

func (f *FuncDelChecker) buildLintError(needCheckInfo *CheckPointerPosition, filePath string, pos *token.Position) *LintError {
	// pos 是需要检测的语句的位置
	lintErr := &LintError{
		Message: NPEMessageTipInfo,
		File:    filePath,
		Line:    pos.Line,
		Colum:   pos.Column,
	}

	if needCheckInfo == nil { // 需要检测但没有检测
		return lintErr
	}

	if needCheckInfo.IsChecked == false {
		return lintErr
	}

	// 虽然检测了， 但引用在先; 这个地方后面看看闭包的地方，要不要细化一下。
	if needCheckInfo.Line > pos.Line || (needCheckInfo.Line == pos.Line && needCheckInfo.Colum > pos.Column) {
		return lintErr
	}
	return nil
}

func Run(pass *analysis.Pass) (interface{}, error) {
	fset := pass.Fset
	var lintErrorList = make([]*LintError, 0)
	for _, file := range pass.Files {
		for _, decl := range file.Decls {
			switch decl := decl.(type) {
			case *ast.FuncDecl:
				checker := InitFuncDelChecker(pass)
				// 记录来自外部的指针变量
				checker.preRecordNilPointerFromOutside(decl, &lintErrorList)
				// 检测外部指针变量的引用
				checker.detectNilPointerReference(decl, fset, &lintErrorList)
			}
		}
	}
	fmt.Println(lintErrorList)

	return nil, nil
}

type FuncDelChecker struct {
	pass *analysis.Pass

	originFieldMap              map[string]string
	needCheckPointerPositionMap map[string][]*CheckPointerPosition
}

// 相同变量名 [position1, position2, position3]
// 记录到有做if nil判断的时候，从最后往前面找到第一个比自己靠前的position 进行IsCheck true , 更新Line ,Colum
// 检测的时候，也从最后往前找到第一个比自己靠前的position， 看看IsCheck 是否为true

func InitFuncDelChecker(pass *analysis.Pass) *FuncDelChecker {
	if pass == nil {
		return nil
	}

	return &FuncDelChecker{
		pass:                        pass,
		originFieldMap:              make(map[string]string),
		needCheckPointerPositionMap: make(map[string][]*CheckPointerPosition),
	}
}

func (f *FuncDelChecker) detectNilPointerReference(
	decl *ast.FuncDecl,
	fset *token.FileSet,
	lintErrorList *[]*LintError) {
	// 针对每个函数内部的变量引用，检测是否存在NPE 问题
	for _, stmt := range decl.Body.List {
		f.detectNPEInStatement(stmt, fset, lintErrorList)
	}
}

func (f *FuncDelChecker) detectNPEInStatement(stmt ast.Stmt, fset *token.FileSet, npeLintErrorListPtr *[]*LintError) {
	switch s := stmt.(type) {
	// if语句块内变量检测
	case *ast.IfStmt:
		f.detectIfStatementBlock(s, fset, npeLintErrorListPtr)

	// 赋值语句变量检测, 如temp := a.b.c
	case *ast.AssignStmt:
		// 右边的表达式里面获取 SelectorExpr list
		f.detectAssignmentStatementBlock(s, fset, npeLintErrorListPtr)

	// 表达式语句，如 fmt.Println(a.b, c)
	case *ast.ExprStmt:
		f.detectExprStatementBlock(s, fset, npeLintErrorListPtr)

	// for循环 例如 for k, v := range GetUserInfoList()
	case *ast.RangeStmt:
		f.detectRangeStatementBlock(s, fset, npeLintErrorListPtr)

	// switch case
	case *ast.SwitchStmt:
		f.detectSwitchStatementBlock(s, fset, npeLintErrorListPtr)

		// defer -- 优先级不高
	}
}

func (f *FuncDelChecker) detectSwitchStatementBlock(s *ast.SwitchStmt, fset *token.FileSet, npeLintErrorListPtr *[]*LintError) {
	if s == nil {
		return
	}

	if s.Body == nil {
		return
	}

	for _, b := range s.Body.List {
		switch bStmt := b.(type) {
		case *ast.CaseClause:
			for _, expr := range bStmt.Body {
				switch expr := expr.(type) {
				case *ast.IfStmt, *ast.AssignStmt, *ast.ExprStmt, *ast.RangeStmt, *ast.SwitchStmt:
					f.detectNPEInStatement(expr, fset, npeLintErrorListPtr)
				}
			}
		}
	}
}

func (f *FuncDelChecker) detectRangeStatementBlock(s *ast.RangeStmt, fset *token.FileSet, npeLintErrorListPtr *[]*LintError) {
	if s == nil {
		return
	}

	if s.Body == nil {
		return
	}

	for _, b := range s.Body.List {
		switch bStmt := b.(type) {
		case *ast.IfStmt, *ast.AssignStmt, *ast.ExprStmt, *ast.RangeStmt, *ast.SwitchStmt:
			f.detectNPEInStatement(bStmt, fset, npeLintErrorListPtr)
		}
	}
}

func (f *FuncDelChecker) detectExprStatementBlock(s *ast.ExprStmt, fset *token.FileSet, lintErrorListPtr *[]*LintError) {
	switch callExpr := s.X.(type) {
	case *ast.CallExpr:
		if selectExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
			lintErrors := f.detectSelectorReferenceWithFunc(selectExpr, fset)
			if len(lintErrors) > 0 {
				*lintErrorListPtr = append(*lintErrorListPtr, lintErrors...)
			}
		}

		for _, expr := range callExpr.Args {
			switch expr := expr.(type) {
			case *ast.SelectorExpr:
				lintErrors := f.detectSelectorReferenceWithFunc(expr, fset)
				if len(lintErrors) > 0 {
					*lintErrorListPtr = append(*lintErrorListPtr, lintErrors...)
				}

			case *ast.CallExpr:
				if selectExpr, ok := expr.Fun.(*ast.SelectorExpr); ok {
					lintErrors := f.detectSelectorReferenceWithFunc(selectExpr, fset)
					if len(lintErrors) > 0 {
						*lintErrorListPtr = append(*lintErrorListPtr, lintErrors...)
					}
				}
			}
		}

	}
}

func (f *FuncDelChecker) detectAssignmentStatementBlock(
	s *ast.AssignStmt,
	fset *token.FileSet,
	npeLintErrorListPtr *[]*LintError) {
	for _, expr := range s.Rhs {
		switch EX := expr.(type) {
		case *ast.SelectorExpr: // 例如 tempC := as.B
			lintErrors := f.detectSelectorReferenceWithFunc(EX, fset)
			if len(lintErrors) > 0 {
				*npeLintErrorListPtr = append(*npeLintErrorListPtr, lintErrors...)
			}

		case *ast.CallExpr:
			if selectExpr, ok := EX.Fun.(*ast.SelectorExpr); ok {
				lintErrors := f.detectSelectorReferenceWithFunc(selectExpr, fset)
				if len(lintErrors) > 0 {
					*npeLintErrorListPtr = append(*npeLintErrorListPtr, lintErrors...)
				}
			}
		}
	}
}

func (f *FuncDelChecker) detectIfStatementBlock(s *ast.IfStmt, fset *token.FileSet, lintErrorListPtr *[]*LintError) {
	for _, stmt := range s.Body.List {
		switch expr := stmt.(type) {
		case *ast.ExprStmt:
			x := expr.X
			switch x := x.(type) {
			case *ast.CallExpr:
				for _, arg := range x.Args {
					switch EX := arg.(type) {
					case *ast.SelectorExpr: // 例如 fmt.Println(a.b.c)
						lintErrors := f.detectSelectorReferenceWithFunc(EX, fset)
						*lintErrorListPtr = append(*lintErrorListPtr, lintErrors...)
					}
				}
			}

		case *ast.IfStmt, *ast.RangeStmt, *ast.SwitchStmt, *ast.AssignStmt: // 有语句块展开，需要递归处理:
			f.detectNPEInStatement(expr, fset, lintErrorListPtr)
		}
	}
}

// 检测赋值语句v:= a.b.c ,表达式语句 fmt.Println(a.b.c) 是否存在NPE
func (f *FuncDelChecker) detectSelectorReferenceWithFunc(
	EX *ast.SelectorExpr,
	fset *token.FileSet,
) []*LintError {
	var (
		nodeList      = make([]*SelectNode, 0)
		lintErrorList []*LintError
	)
	walkSelectorWithFunc(EX, fset, func(node ast.Node) { // 递归获取变量名
		switch exprInner := node.(type) {
		case *ast.SelectorExpr:
			if exprInner != nil && exprInner.Sel != nil {
				nodeType, isReturnSingleFunc := GetNodeType(exprInner.Sel, f.pass.TypesInfo)
				nodeList = append(nodeList, &SelectNode{
					Name:               exprInner.Sel.Name,
					CurIdent:           exprInner.Sel,
					Type:               nodeType,
					IsReturnSingleFunc: isReturnSingleFunc})
			}

		case *ast.CallExpr:
			switch exFunc := exprInner.Fun.(type) {
			case *ast.SelectorExpr:
				if exFunc != nil && exFunc.Sel != nil {
					nodeType, isReturnSingleFunc := GetNodeType(exFunc.Sel, f.pass.TypesInfo)
					nodeList = append(nodeList, &SelectNode{
						Name:               exFunc.Sel.Name,
						CurIdent:           exFunc.Sel,
						Type:               nodeType,
						IsReturnSingleFunc: isReturnSingleFunc})
				}

			case *ast.Ident:
				nodeType, isReturnSingleFunc := GetNodeType(exFunc, f.pass.TypesInfo)
				nodeList = append(nodeList, &SelectNode{
					Name:               exFunc.Name,
					CurIdent:           exFunc,
					Type:               nodeType,
					IsReturnSingleFunc: isReturnSingleFunc})
			}

		case *ast.Ident: // 最后一个，会结束递归
			if exprInner != nil {
				nodeType, isReturnSingleFunc := GetNodeType(exprInner, f.pass.TypesInfo)
				nodeList = append(nodeList, &SelectNode{
					Name:               exprInner.Name,
					CurIdent:           exprInner,
					Type:               nodeType,
					IsReturnSingleFunc: isReturnSingleFunc})
				sort.SliceStable(nodeList, func(i, j int) bool { // 反转数组
					return i > j
				})

				var pos token.Position
				GetIdentPosition(&pos, exprInner, fset)
				// 检测a -> a.b -> a.b.c
				//lintErrorList = f.hasSequenceDetectNodeName(fieldName, pos, exprInner)
				fieldName := buildFieldNameFromNodes(nodeList)
				//_, funcIndex := f.isExistFuncRetPtr(nodeList, true)
				isRootFromOutside := f.isRootComeFromOutside(fieldName)
				lintErrorList = f.hasSequenceDetectNode(nodeList, pos, exprInner, isRootFromOutside)
			}
		}
	})
	return lintErrorList
}

// 检测赋值语句v:= a.b.c 右侧的赋值是否存在NPE
//func (f *FuncDelChecker) detectSelectorReference(
//	EX *ast.SelectorExpr,
//	fset *token.FileSet,
//) []*LintError {
//	var (
//		nodeNameList  = make([]*SelectNode, 0)
//		lintErrorList []*LintError
//	)
//	WalkSelector(EX, fset, func(node ast.Node) { // 递归获取变量名
//		switch exprInner := node.(type) {
//		case *ast.SelectorExpr:
//			if exprInner != nil && exprInner.Sel != nil {
//				nodeType := GetNodeType(exprInner.Sel, f.pass.TypesInfo)
//				nodeNameList = append(nodeNameList, &SelectNode{Name: exprInner.Sel.Name, Type: nodeType})
//			}
//
//		case *ast.Ident: // 最后一个，会结束递归
//			if exprInner != nil {
//				nodeType := GetNodeType(exprInner, f.pass.TypesInfo)
//				nodeNameList = append(nodeNameList, &SelectNode{Name: exprInner.Name, Type: nodeType})
//				sort.SliceStable(nodeNameList, func(i, j int) bool { // 反转数组
//					return i > j
//				})
//				//var fieldName = strings.Join(nodeNameList, ".") // 如 a.b.c
//				//var fieldName = buildFieldNameFromNodes(nodeNameList)
//				var pos token.Position
//				GetIdentPosition(&pos, exprInner, fset)
//				//fmt.Printf("表达式这 fieldName:%s, Line: %d, Column: %d\n", fieldName, pos.Line, pos.Column)
//				// 检测a -> a.b -> a.b.c
//				//lintErrorList = f.hasSequenceDetectNodeName(fieldName, pos, exprInner)
//				lintErrorList = f.hasSequenceDetectNode(nodeNameList, pos, exprInner)
//				//lintError := getPotentialNilPointerReference(fieldName, pos.Filename, &pos, originFieldMap, needCheckPointerPositionMap)
//
//			}
//		}
//	})
//	return lintErrorList
//}

// 正常流程得依次检测a, a.b, a.b.c , 所以检测a.b.c 之前，看看前面的是否已经校验了
func (f *FuncDelChecker) hasSequenceDetectNode(
	nodeList []*SelectNode,
	pos token.Position,
	expr *ast.Ident,
	// funcIndex int, //o.GetUserInfo().GetMembership() -- funcIndex: 1
	isRootFromOutside bool,
) []*LintError {
	var result = make([]*LintError, 0)
	if len(nodeList) <= 1 {
		return nil
	}

	for len(nodeList) > 0 { // 还是只检测root node from outside 稳一点, 避免误伤
		//if isRootFromOutside == false && funcIndex >= 0 && len(nodeList)-1 <= funcIndex {
		//	return result // root 不是来自外部 ,funcIndex 之前的不用检测
		//}

		if len(nodeList) >= 2 && nodeList[len(nodeList)-2].Type == NodeTypeNonSinglePtr {
			nodeList = nodeList[0 : len(nodeList)-1]
			continue
		}

		fieldName := buildFieldNameFromNodes(nodeList)
		if len(nodeList) >= 2 && nodeList[len(nodeList)-2].CurIdent != nil {
			expr = nodeList[len(nodeList)-2].CurIdent
		}

		lintError := f.getPotentialNilPointerReference(fieldName, pos.Filename, &pos, expr, isRootFromOutside)
		if lintError != nil {
			result = append(result, lintError)
		}
		nodeList = nodeList[0 : len(nodeList)-1]
	}

	return result
}

// 正常流程得依次检测a, a.b, a.b.c , 所以检测a.b.c 之前，看看前面的是否已经校验了
//func (f *FuncDelChecker) hasSequenceDetectNodeName(
//	fieldName string,
//	pos token.Position,
//	expr *ast.Ident,
//) []*LintError {
//	var result = make([]*LintError, 0)
//	nodeNameList := strings.Split(fieldName, ".")
//	if len(nodeNameList) <= 1 {
//		return nil
//	}
//
//	for len(nodeNameList) > 0 {
//		fieldName = strings.Join(nodeNameList, ".")
//		lintError := f.getPotentialNilPointerReference(fieldName, pos.Filename, &pos, expr)
//		if lintError != nil {
//			result = append(result, lintError)
//		}
//		nodeNameList = nodeNameList[0 : len(nodeNameList)-1]
//	}
//
//	return result
//}

func GetFuncSignature(ex *ast.CallExpr, typeInfo *types.Info) *types.Signature {
	if ex == nil {
		return nil
	}

	var (
		sig *types.Signature
		ok  bool
	)

	fn := ex.Fun
	switch fn := fn.(type) { // 需要找到叶子节点
	case *ast.Ident:
		obj := typeInfo.ObjectOf(fn)
		sig, ok = obj.Type().(*types.Signature)
		if !ok {
			return nil
		}

	case *ast.SelectorExpr:
		obj := typeInfo.ObjectOf(fn.Sel)
		sig, ok = obj.Type().(*types.Signature)
		if !ok {
			return nil
		}
	}

	return sig
}
