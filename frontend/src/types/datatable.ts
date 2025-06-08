// Table 資料結構的類型定義
export interface Column {
  id: number;
  name: string;
}

export interface Cell {
  [key: string]: string | number;
}

export interface Row {
  id: number;
  cells: Cell;
}

export interface TableData {
  columns: Column[];
  rows: Row[];
}

// 儲存格編輯狀態
export interface EditingState {
  tableName: string;
  rowIndex: number;
  colIndex: number;
  colName: string;
  value: string;
  isEditing: boolean;
}

// 基於ID的儲存格編輯狀態
export interface EditingStateByID {
  tableID: number;
  rowIndex: number;
  colIndex: number;
  colName: string;
  value: string;
  isEditing: boolean;
}

// 表格基本信息
export interface TableInfo {
  id: number;
  name: string;
  rowCount: number;
  colCount: number;
}
