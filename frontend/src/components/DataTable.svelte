<script lang="ts">
  import { onMount, createEventDispatcher } from "svelte";
  import { TableData, EditingStateByID } from "../types/datatable";
  import {
    GetTableDataByID,
    UpdateCellValueByID,
    UpdateColumnNameByID,
    AddRowByID,
    AddColumnByID,
    GetText,
  } from "../../wailsjs/go/main/App";
  import ContextMenu from "./ContextMenu.svelte";
  import type { ContextMenuConfig } from "../types/contextMenu";

  // 組件屬性
  export let tableID: number;
  export let tableKey: number = 0; // 用於強制重新載入的 key

  // 新增：表格縮放比例
  let tableScale = 1;

  // 創建事件分發器
  const dispatch = createEventDispatcher();

  // i18n 翻譯輔助函數
  async function t(key: string, vars?: Record<string, any>): Promise<string> {
    try {
      let text = await GetText(key);
      if (vars) {
        Object.entries(vars).forEach(([key, value]) => {
          text = text.replace(`{${key}}`, value);
        });
      }
      return text;
    } catch (error) {
      console.warn(`Translation missing for key: ${key}`);
      return key;
    }
  }
  // 翻譯文字快取
  let texts: Record<string, string> = {};

  // 載入翻譯文字
  async function loadTexts() {
    const keys = [
      "ui.table.loading",
      "ui.table.no_data",
      "ui.table.selected_content",
      "ui.table.selected_row",
      "ui.table.selected_column",
      "ui.table.cell_position",
      "ui.table.update_failed",
      "ui.table.zoom", // 新增翻譯鍵
      "ui.context_menu.insert_row_above",
      "ui.context_menu.insert_row_below",
      "ui.context_menu.duplicate_row",
      "ui.context_menu.delete_row",
      "ui.context_menu.insert_column_left",
      "ui.context_menu.insert_column_right",
      "ui.context_menu.rename_column",
      "ui.context_menu.duplicate_column",
      "ui.context_menu.delete_column",
      "ui.context_menu.copy",
      "ui.context_menu.paste",
      "ui.context_menu.clear",
    ];

    for (const key of keys) {
      try {
        texts[key] = await GetText(key);
      } catch (error) {
        console.warn(`Failed to load translation for ${key}`);
      }
    }
  }
  // 狀態變數
  let tableData: TableData | null = null;
  let loading = true;
  let error = "";
  let lastTableID = -1;
  let lastTableKey = -1;

  // 編輯狀態
  let editingState: EditingStateByID = {
    tableID: -1,
    rowIndex: -1,
    colIndex: -1,
    colName: "",
    value: "",
    isEditing: false,
  }; // 選中狀態
  let selectedRow = -1;
  let selectedCol = -1;
  let selectedCellContent = "";

  // 選擇模式：'cell' | 'row' | 'column' | 'range'
  let selectionMode = "cell";
  let selectedRowRange = new Set(); // 選中的行範圍
  let selectedColRange = new Set(); // 選中的列範圍

  // 範圍選取狀態
  let rangeSelectStartRow = -1;
  let rangeSelectStartCol = -1;
  let rangeSelectEndRow = -1;
  let rangeSelectEndCol = -1;
  let isSelectingRange = false;
  let isDragging = false;

  // 剪貼簿資料
  let clipboardData: string[][] = [];
  let clipboardType: "copy" | "cut" | null = null;

  // 右鍵菜單狀態
  let contextMenuVisible = false;
  let contextMenuX = 0;
  let contextMenuY = 0;
  let contextMenuType = ""; // 'row' | 'column' | 'cell'
  let contextMenuContext = {}; // 上下文信息  // 右鍵菜單配置
  let contextMenuConfig: ContextMenuConfig = {
    row: [
      { id: "insertRowAbove", label: "", icon: "⬆️" },
      { id: "insertRowBelow", label: "", icon: "⬇️" },
      { id: "separator1", type: "separator" },
      { id: "duplicateRow", label: "", icon: "📋" },
      { id: "paste", label: "", icon: "📄", disabled: true },
      { id: "separator2", type: "separator" },
      { id: "deleteRow", label: "", icon: "🗑️", danger: true },
    ],
    column: [
      { id: "insertColumnLeft", label: "", icon: "⬅️" },
      { id: "insertColumnRight", label: "", icon: "➡️" },
      { id: "separator1", type: "separator" },
      { id: "renameColumn", label: "", icon: "✏️" },
      { id: "duplicateColumn", label: "", icon: "📋" },
      { id: "paste", label: "", icon: "📄", disabled: true },
      { id: "separator2", type: "separator" },
      { id: "deleteColumn", label: "", icon: "🗑️", danger: true },
    ],
    cell: [
      { id: "copy", label: "", icon: "📋" },
      { id: "paste", label: "", icon: "📄", disabled: true },
      { id: "separator1", type: "separator" },
      { id: "clear", label: "", icon: "🧹" },
      { id: "separator2", type: "separator" },
      { id: "insertRowAbove", label: "", icon: "⬆️" },
      { id: "insertRowBelow", label: "", icon: "⬇️" },
      { id: "insertColumnLeft", label: "", icon: "⬅️" },
      { id: "insertColumnRight", label: "", icon: "➡️" },
    ],
    range: [
      { id: "copy", label: "", icon: "📋" },
      { id: "paste", label: "", icon: "📄", disabled: true },
      { id: "separator1", type: "separator" },
      { id: "clear", label: "", icon: "🧹" },
      { id: "separator2", type: "separator" },
      { id: "fillSeries", label: "", icon: "📊" },
    ],
  };
  // 更新菜單配置的翻譯文字
  function updateMenuLabels() {
    // Row menu
    contextMenuConfig.row[0].label =
      texts["ui.context_menu.insert_row_above"] || "在上方插入列";
    contextMenuConfig.row[1].label =
      texts["ui.context_menu.insert_row_below"] || "在下方插入列";
    contextMenuConfig.row[3].label =
      texts["ui.context_menu.duplicate_row"] || "複製列";
    contextMenuConfig.row[4].label = texts["ui.context_menu.paste"] || "貼上";
    contextMenuConfig.row[4].disabled = clipboardData.length === 0; // 動態禁用貼上
    contextMenuConfig.row[6].label =
      texts["ui.context_menu.delete_row"] || "刪除列";

    // Column menu
    contextMenuConfig.column[0].label =
      texts["ui.context_menu.insert_column_left"] || "在左邊插入變項";
    contextMenuConfig.column[1].label =
      texts["ui.context_menu.insert_column_right"] || "在右邊插入變項";
    contextMenuConfig.column[3].label =
      texts["ui.context_menu.rename_column"] || "重新命名變項";
    contextMenuConfig.column[4].label =
      texts["ui.context_menu.duplicate_column"] || "複製變項";
    contextMenuConfig.column[5].label =
      texts["ui.context_menu.paste"] || "貼上";
    contextMenuConfig.column[5].disabled = clipboardData.length === 0; // 動態禁用貼上
    contextMenuConfig.column[7].label =
      texts["ui.context_menu.delete_column"] || "刪除變項"; // Cell menu - 動態更新貼上選項狀態
    contextMenuConfig.cell[0].label = texts["ui.context_menu.copy"] || "複製";
    contextMenuConfig.cell[1].label = texts["ui.context_menu.paste"] || "貼上";
    contextMenuConfig.cell[1].disabled = clipboardData.length === 0; // 動態禁用貼上
    contextMenuConfig.cell[3].label =
      texts["ui.context_menu.clear"] || "清除內容";
    contextMenuConfig.cell[5].label =
      texts["ui.context_menu.insert_row_above"] || "在上方插入列";
    contextMenuConfig.cell[6].label =
      texts["ui.context_menu.insert_row_below"] || "在下方插入列";
    contextMenuConfig.cell[7].label =
      texts["ui.context_menu.insert_column_left"] || "在左邊插入變項";
    contextMenuConfig.cell[8].label =
      texts["ui.context_menu.insert_column_right"] || "在右邊插入變項";

    // Range menu - 範圍選取菜單
    contextMenuConfig.range[0].label = texts["ui.context_menu.copy"] || "複製";
    contextMenuConfig.range[1].label = texts["ui.context_menu.paste"] || "貼上";
    contextMenuConfig.range[1].disabled = clipboardData.length === 0; // 動態禁用貼上
    contextMenuConfig.range[3].label =
      texts["ui.context_menu.clear"] || "清除內容";
    contextMenuConfig.range[5].label = "填充數列"; // 新功能，暫時硬編碼
  }

  // 防止雙擊時觸發點擊的標記
  let doubleClickInProgress = false;
  // 編輯輸入元素引用
  let editInput: HTMLInputElement; // 當進入編輯模式時，設置焦點
  $: if (editingState.isEditing && editInput) {
    setTimeout(() => {
      editInput.focus(); // 僅聚焦，不選取文字
      // editInput.select(); // 暫時移除此行
    }, 0);
  } // 響應式更新選中內容顯示
  $: if (tableData && !editingState.isEditing) {
    updateSelectedCellContent();
  }

  // 當選擇狀態變化時即時更新選中內容顯示
  $: if (tableData && (selectedRow >= 0 || selectedCol >= 0 || selectionMode)) {
    if (!editingState.isEditing) {
      updateSelectedCellContent();
    }
  }
  onMount(async () => {
    lastTableID = tableID;
    lastTableKey = tableKey;

    // 載入翻譯文字
    await loadTexts();
    updateMenuLabels();

    await loadTableData(); // 添加文檔點擊事件監聽器
    document.addEventListener("click", handleDocumentClick);
    document.addEventListener("mouseup", handleGlobalMouseUp);

    // 添加鍵盤事件監聽器
    document.addEventListener("keydown", handleGlobalKeyDown);
    document.addEventListener("keyup", handleGlobalKeyUp);

    return () => {
      // 清理事件監聽器
      document.removeEventListener("click", handleDocumentClick);
      document.removeEventListener("mouseup", handleGlobalMouseUp);
      document.removeEventListener("keydown", handleGlobalKeyDown);
      document.removeEventListener("keyup", handleGlobalKeyUp);
    };
  });

  // 更新選中內容顯示的函數
  function updateSelectedCellContent() {
    if (!tableData) return;

    if (
      selectionMode === "range" &&
      rangeSelectStartRow >= 0 &&
      rangeSelectStartCol >= 0 &&
      rangeSelectEndRow >= 0 &&
      rangeSelectEndCol >= 0
    ) {
      const startRow = Math.min(rangeSelectStartRow, rangeSelectEndRow);
      const endRow = Math.max(rangeSelectStartRow, rangeSelectEndRow);
      const startCol = Math.min(rangeSelectStartCol, rangeSelectEndCol);
      const endCol = Math.max(rangeSelectStartCol, rangeSelectEndCol);
      const rowCount = endRow - startRow + 1;
      const colCount = endCol - startCol + 1;
      selectedCellContent = `已選取 ${rowCount} 行 × ${colCount} 欄 (${rowCount * colCount} 個儲存格)`;
    } else if (selectionMode === "row" && selectedRow >= 0) {
      selectedCellContent = (
        texts["ui.table.selected_row"] || "第 {row} 列"
      ).replace("{row}", (selectedRow + 1).toString());
    } else if (selectionMode === "column" && selectedCol >= 0) {
      selectedCellContent = (
        texts["ui.table.selected_column"] || "{column} 變項"
      ).replace("{column}", indexToLetters(selectedCol));
    } else if (
      selectionMode === "cell" &&
      selectedRow >= 0 &&
      selectedCol >= 0
    ) {
      const column = tableData.columns[selectedCol];
      if (column) {
        const cellValue = tableData.rows[selectedRow]?.cells[column.name];
        const displayValue = formatCellValue(cellValue);
        const position = (texts["ui.table.cell_position"] || "{column}{row}")
          .replace("{column}", indexToLetters(selectedCol))
          .replace("{row}", (selectedRow + 1).toString());
        selectedCellContent = `${position}: ${displayValue}`;
      }
    } else {
      selectedCellContent = "";
    }
  }

  // 全域鍵盤事件處理
  function handleGlobalKeyDown(event: KeyboardEvent) {
    // 如果在編輯狀態，不處理快捷鍵
    if (editingState.isEditing) return;

    // Ctrl/Cmd + C 複製
    if ((event.ctrlKey || event.metaKey) && event.key === "c") {
      event.preventDefault();
      handleCopy();
    }
    // Ctrl/Cmd + V 貼上
    else if ((event.ctrlKey || event.metaKey) && event.key === "v") {
      event.preventDefault();
      handlePaste();
    }
    // Escape 清除選取
    else if (event.key === "Escape") {
      clearSelection();
    }
    // Shift 按下開始範圍選取模式
    else if (event.key === "Shift" && !isSelectingRange) {
      if (selectedRow >= 0 && selectedCol >= 0) {
        startRangeSelection(selectedRow, selectedCol);
      }
    }
  }

  function handleGlobalKeyUp(event: KeyboardEvent) {
    // Shift 放開結束範圍選取模式
    if (event.key === "Shift" && isSelectingRange) {
      endRangeSelection();
    }
  }
  // 開始範圍選取
  function startRangeSelection(row: number, col: number) {
    isSelectingRange = true;
    rangeSelectStartRow = row;
    rangeSelectStartCol = col;
    rangeSelectEndRow = row;
    rangeSelectEndCol = col;
    selectionMode = "range";

    // 清除單格和行列選取狀態
    selectedRow = -1;
    selectedCol = -1;
    selectedRowRange = new Set();
    selectedColRange = new Set();
  }

  // 結束範圍選取
  function endRangeSelection() {
    isSelectingRange = false;
  }

  // 更新範圍選取
  function updateRangeSelection(row: number, col: number) {
    if (isSelectingRange) {
      rangeSelectEndRow = row;
      rangeSelectEndCol = col;
      selectionMode = "range";
    }
  }

  // 清除選取
  function clearSelection() {
    selectionMode = "cell";
    selectedRow = -1;
    selectedCol = -1;
    selectedRowRange = new Set();
    selectedColRange = new Set();
    rangeSelectStartRow = -1;
    rangeSelectStartCol = -1;
    rangeSelectEndRow = -1;
    rangeSelectEndCol = -1;
    isSelectingRange = false;
  }

  // 複製功能
  function handleCopy() {
    if (!tableData) return;

    let dataToCopy: string[][] = [];
    if (
      selectionMode === "range" &&
      rangeSelectStartRow >= 0 &&
      rangeSelectEndRow >= 0
    ) {
      // 複製選取範圍
      const startRow = Math.min(rangeSelectStartRow, rangeSelectEndRow);
      const endRow = Math.max(rangeSelectStartRow, rangeSelectEndRow);
      const startCol = Math.min(rangeSelectStartCol, rangeSelectEndCol);
      const endCol = Math.max(rangeSelectStartCol, rangeSelectEndCol);

      for (let row = startRow; row <= endRow; row++) {
        const rowData: string[] = [];
        for (let col = startCol; col <= endCol; col++) {
          const column = tableData.columns[col];
          if (column && tableData.rows[row]) {
            const cellValue = tableData.rows[row].cells[column.name];
            rowData.push(formatCellValue(cellValue));
          } else {
            rowData.push("");
          }
        }
        dataToCopy.push(rowData);
      }
    } else if (
      selectionMode === "cell" &&
      selectedRow >= 0 &&
      selectedCol >= 0
    ) {
      // 複製單一儲存格
      const column = tableData.columns[selectedCol];
      if (column && tableData.rows[selectedRow]) {
        const cellValue = tableData.rows[selectedRow].cells[column.name];
        dataToCopy = [[formatCellValue(cellValue)]];
      }
    } else if (selectionMode === "row" && selectedRow >= 0) {
      // 複製整行
      const rowData: string[] = [];
      for (let col = 0; col < tableData.columns.length; col++) {
        const column = tableData.columns[col];
        if (column && tableData.rows[selectedRow]) {
          const cellValue = tableData.rows[selectedRow].cells[column.name];
          rowData.push(formatCellValue(cellValue));
        } else {
          rowData.push("");
        }
      }
      dataToCopy = [rowData];
    } else if (selectionMode === "column" && selectedCol >= 0) {
      // 複製整欄
      const column = tableData.columns[selectedCol];
      if (column) {
        for (let row = 0; row < tableData.rows.length; row++) {
          if (tableData.rows[row]) {
            const cellValue = tableData.rows[row].cells[column.name];
            dataToCopy.push([formatCellValue(cellValue)]);
          } else {
            dataToCopy.push([""]);
          }
        }
      }
    }

    if (dataToCopy.length > 0) {
      clipboardData = dataToCopy;
      clipboardType = "copy";

      // 複製到系統剪貼簿
      const textData = dataToCopy.map((row) => row.join("\t")).join("\n");
      navigator.clipboard.writeText(textData).catch((err) => {
        console.warn("無法複製到系統剪貼簿:", err);
      });

      console.log("已複製資料:", dataToCopy);
    }
  }
  // 貼上功能
  async function handlePaste() {
    if (!tableData || clipboardData.length === 0) return;

    try {
      if (editingState.isEditing) {
        // 在編輯狀態時，將所有內容插入同一格
        const allText = clipboardData.map((row) => row.join(" ")).join(" ");
        editingState.value += allText;
      } else {
        // 不在編輯狀態時，自動溢出貼上
        let startRow = selectedRow >= 0 ? selectedRow : 0;
        let startCol = selectedCol >= 0 ? selectedCol : 0;

        // 計算需要的最大行數和列數
        const requiredRows = startRow + clipboardData.length;
        const requiredCols =
          startCol + Math.max(...clipboardData.map((row) => row.length)); // 檢查是否需要添加列
        const currentColCount = tableData.columns.length;
        if (requiredCols > currentColCount) {
          const colsToAdd = requiredCols - currentColCount;
          console.log(`需要添加 ${colsToAdd} 個列`);
          for (let i = 0; i < colsToAdd; i++) {
            const newColName = ""; // 自動擴張的欄不要有名字
            const success = await AddColumnByID(tableID, newColName);
            if (!success) {
              console.error(`添加列失敗`);
              break;
            }
          }

          // 重新載入資料以獲取新的列結構
          await loadTableData();
          // 確保 tableData 已更新
          if (!tableData) {
            error = "擴張表格後無法載入資料";
            return;
          }
        }

        // 檢查是否需要添加行
        const currentRowCount = tableData.rows.length;
        if (requiredRows > currentRowCount) {
          const rowsToAdd = requiredRows - currentRowCount;
          console.log(`需要添加 ${rowsToAdd} 個行`);

          for (let i = 0; i < rowsToAdd; i++) {
            const success = await AddRowByID(tableID);
            if (!success) {
              console.error(`添加行失敗`);
              break;
            }
          }

          // 重新載入資料以獲取新的行結構
          await loadTableData();
          // 確保 tableData 已更新
          if (!tableData) {
            error = "擴張表格後無法載入資料";
            return;
          }
        }

        // 現在執行貼上操作
        for (let rowOffset = 0; rowOffset < clipboardData.length; rowOffset++) {
          const targetRow = startRow + rowOffset;
          // 確保目標行存在（應該在上面的擴張中已經處理）
          if (targetRow >= tableData.rows.length) continue;

          const rowData = clipboardData[rowOffset];
          for (let colOffset = 0; colOffset < rowData.length; colOffset++) {
            const targetCol = startCol + colOffset;
            // 確保目標列存在（應該在上面的擴張中已經處理）
            if (targetCol >= tableData.columns.length) continue;

            const column = tableData.columns[targetCol];
            if (column) {
              const processedValue = parseInputValue(rowData[colOffset]);
              await UpdateCellValueByID(
                tableID,
                targetRow,
                targetCol,
                processedValue
              );
            }
          }
        }

        // 重新載入資料
        await loadTableData();
      }
    } catch (err) {
      error = `貼上失敗: ${err}`;
    }
  }

  // 貼上到整列功能
  async function handlePasteToRow(rowIndex: number) {
    if (!tableData || clipboardData.length === 0) return;

    try {
      // 獲取第一行資料來貼上到指定列
      const firstRowData = clipboardData[0];

      // 確保目標列有足夠的欄位
      const requiredCols = firstRowData.length;
      const currentColCount = tableData.columns.length;
      if (requiredCols > currentColCount) {
        const colsToAdd = requiredCols - currentColCount;
        for (let i = 0; i < colsToAdd; i++) {
          const newColName = ""; // 自動擴張的欄不要有名字
          const success = await AddColumnByID(tableID, newColName);
          if (!success) {
            console.error(`添加列失敗`);
            break;
          }
        }
        // 重新載入資料
        await loadTableData();
        if (!tableData) return;
      }

      // 貼上資料到指定列
      for (let colIndex = 0; colIndex < firstRowData.length; colIndex++) {
        if (colIndex >= tableData.columns.length) break; // Corrected syntax
        const column = tableData.columns[colIndex];
        if (column) {
          const processedValue = parseInputValue(firstRowData[colIndex]);
          await UpdateCellValueByID(
            tableID,
            rowIndex,
            colIndex,
            processedValue
          );
        }
      }

      // 重新載入資料
      await loadTableData();
    } catch (err) {
      error = `貼上到列失敗: ${err}`;
    }
  }

  // 貼上到整欄功能
  async function handlePasteToColumn(colIndex: number) {
    if (!tableData || clipboardData.length === 0) return;

    try {
      // 確保目標欄位存在
      if (colIndex >= tableData.columns.length) return;

      // 確保有足夠的列
      const requiredRows = clipboardData.length;
      const currentRowCount = tableData.rows.length;
      if (requiredRows > currentRowCount) {
        const rowsToAdd = requiredRows - currentRowCount;
        for (let i = 0; i < rowsToAdd; i++) {
          const success = await AddRowByID(tableID);
          if (!success) {
            console.error(`添加行失敗`);
            break;
          }
        }
        // 重新載入資料
        await loadTableData();
        if (!tableData) return;
      }

      const column = tableData.columns[colIndex];
      if (!column) return;

      // 貼上資料到指定欄位
      for (let rowOffset = 0; rowOffset < clipboardData.length; rowOffset++) {
        if (rowOffset >= tableData.rows.length) break;

        // 取得該行剪貼簿資料的第一個值
        const rowData = clipboardData[rowOffset];
        const valueToSet = rowData.length > 0 ? rowData[0] : "";

        const processedValue = parseInputValue(valueToSet);
        await UpdateCellValueByID(tableID, rowOffset, colIndex, processedValue);
      }

      // 重新載入資料
      await loadTableData();
    } catch (err) {
      error = `貼上到欄失敗: ${err}`;
    }
  }

  onMount(async () => {
    lastTableID = tableID;
    lastTableKey = tableKey;

    // 載入翻譯文字
    await loadTexts();
    updateMenuLabels();

    await loadTableData();

    // 添加文檔點擊事件監聽器
    document.addEventListener("click", handleDocumentClick);

    return () => {
      // 清理事件監聽器
      document.removeEventListener("click", handleDocumentClick);
    };
  });

  // 當 tableID 或 tableKey 變化時重新載入
  $: if (
    (tableID !== lastTableID && tableID >= 0) ||
    (tableKey !== lastTableKey && tableKey >= 0)
  ) {
    lastTableID = tableID;
    lastTableKey = tableKey;
    loadTableData();
  } // 載入表格資料
  async function loadTableData() {
    // 檢查 tableID 是否有效
    if (tableID < 0) {
      error = "無效的資料表 ID";
      loading = false;
      return;
    }

    try {
      loading = true;
      error = "";
      const data = await GetTableDataByID(tableID);
      tableData = data as TableData;

      // 計算並分發統計數據
      if (tableData) {
        const stats = calculateStatistics(tableData);
        dispatch("statsUpdate", stats);

        // 設定初始選中狀態到第一個儲存格（如果存在）
        if (
          tableData.rows &&
          tableData.rows.length > 0 &&
          tableData.columns &&
          tableData.columns.length > 0
        ) {
          selectedRow = 0;
          selectedCol = 0;
          selectionMode = "cell";
          selectedRowRange = new Set();
          selectedColRange = new Set();
          // updateSelectedCellContent 會由響應式語句自動調用
        } else {
          // 如果沒有資料，重置選中狀態
          selectedRow = -1;
          selectedCol = -1;
          selectedCellContent = "";
          selectionMode = "cell";
          selectedRowRange = new Set();
          selectedColRange = new Set();
        }
      }
    } catch (err) {
      error = `載入資料表失敗: ${err}`;
      tableData = null;
    } finally {
      loading = false;
    }
  }

  // 計算統計數據
  function calculateStatistics(data: TableData) {
    const rowCount = data.rows ? data.rows.length : 0;
    const colCount = data.columns ? data.columns.length : 0;
    const totalCells = rowCount * colCount;

    // 計算數值欄位數量
    let numericCols = 0;
    if (data.columns && data.rows) {
      data.columns.forEach((col) => {
        let hasNumeric = false;
        // 檢查前10行來判斷是否為數值欄位
        for (let i = 0; i < Math.min(10, rowCount); i++) {
          const row = data.rows[i];
          if (row && row.cells) {
            const value = row.cells[col.name];
            if (value && !isNaN(Number(value))) {
              hasNumeric = true;
              break;
            }
          }
        }
        if (hasNumeric) numericCols++;
      });
    }
    return {
      total_rows: rowCount.toString(),
      total_variables: colCount.toString(),
      total_cells: totalCells.toString(),
      numeric_variables: numericCols.toString(),
    };
  } // 儲存格點擊處理
  function handleCellClick(
    rowIndex: number,
    colIndex: number,
    value: string,
    event?: Event
  ) {
    // 如果正在雙擊過程中，忽略點擊事件
    if (doubleClickInProgress) {
      return;
    } // 檢查是否按住 Shift 鍵進行範圍選取
    if (
      event &&
      "shiftKey" in event &&
      (event as MouseEvent).shiftKey &&
      selectedRow >= 0 &&
      selectedCol >= 0
    ) {
      startRangeSelection(selectedRow, selectedCol);
      updateRangeSelection(rowIndex, colIndex);
      return;
    }

    // 如果正在選取範圍且按住 Shift
    if (isSelectingRange) {
      updateRangeSelection(rowIndex, colIndex);
      return;
    }

    // 如果正在編輯其他儲存格（不是當前點擊的格子），先結束編輯
    if (editingState.isEditing) {
      // 檢查是否點擊的是同一個格子
      const isSameCell =
        editingState.rowIndex === rowIndex &&
        editingState.colIndex === colIndex;

      // 如果是同一個格子，只更新選擇狀態，不更新顯示內容和編輯狀態
      if (isSameCell) {
        selectedRow = rowIndex;
        selectedCol = colIndex;
        // 不更新 selectedCellContent，保持編輯中的內容
        return;
      }

      // 只有在點擊不同格子時才結束編輯
      handleEditComplete();
    }

    // 清除範圍選取狀態
    isSelectingRange = false;
    rangeSelectStartRow = -1;
    rangeSelectStartCol = -1;
    rangeSelectEndRow = -1;
    rangeSelectEndCol = -1;

    // 更新選擇狀態為儲存格模式
    selectionMode = "cell";
    selectedRow = rowIndex;
    selectedCol = colIndex;
    selectedRowRange = new Set();
    selectedColRange = new Set();
    // selectedCellContent 會自動由響應式語句更新
  }
  // 儲存格雙擊處理 (進入編輯模式)
  function handleCellDblClick(
    rowIndex: number,
    colIndex: number,
    colName: string,
    rawValue: any
  ) {
    // 設置雙擊標記
    doubleClickInProgress = true;

    // 如果已經在編輯同一個格子，不要重新進入編輯模式
    if (
      editingState.isEditing &&
      editingState.rowIndex === rowIndex &&
      editingState.colIndex === colIndex
    ) {
      // 清除雙擊標記
      setTimeout(() => {
        doubleClickInProgress = false;
      }, 10);
      return;
    }

    // 如果正在編輯其他格子，先結束編輯
    if (editingState.isEditing) {
      handleEditComplete();
    }

    // 更新選擇狀態，但不更新顯示內容（避免刷新正在輸入的內容）
    selectedRow = rowIndex;
    selectedCol = colIndex;
    // 不更新 selectedCellContent

    // 進入編輯模式 - 對 nil 值使用空字串而不是點號
    const editValue =
      rawValue === null || rawValue === undefined ? "" : String(rawValue);
    editingState = {
      tableID,
      rowIndex,
      colIndex,
      colName,
      value: editValue,
      isEditing: true,
    };

    // 清除雙擊標記
    setTimeout(() => {
      doubleClickInProgress = false;
    }, 10);
  }

  // 鼠標按下處理 - 開始拖拽選取
  function handleCellMouseDown(
    rowIndex: number,
    colIndex: number,
    event: MouseEvent
  ) {
    if (event.button !== 0) return; // 只處理左鍵

    // 如果按住 Shift，開始範圍選取
    if (event.shiftKey && selectedRow >= 0 && selectedCol >= 0) {
      startRangeSelection(selectedRow, selectedCol);
      updateRangeSelection(rowIndex, colIndex);
      return;
    }

    // 開始拖拽選取
    isDragging = true;
    startRangeSelection(rowIndex, colIndex);

    // 防止文字選取
    event.preventDefault();
  }

  // 鼠標放開處理
  function handleCellMouseUp() {
    if (isDragging) {
      isDragging = false;
      endRangeSelection();
    }
  }

  // 鼠標進入處理
  function handleCellMouseEnter(rowIndex: number, colIndex: number) {
    if (isDragging || isSelectingRange) {
      updateRangeSelection(rowIndex, colIndex);
    }
  }

  // 欄位標題點擊處理
  function handleColumnHeaderClick(colIndex: number, colName: string) {
    // 如果正在雙擊過程中，忽略點擊事件
    if (doubleClickInProgress) {
      return;
    }

    // 如果正在編輯其他欄位標題（不是當前點擊的欄位），先結束編輯
    if (editingState.isEditing) {
      // 檢查是否點擊的是同一個欄位標題
      const isSameHeader =
        editingState.rowIndex === -1 && editingState.colIndex === colIndex;

      // 如果是同一個欄位標題，只更新選擇狀態，不更新顯示內容和編輯狀態
      if (isSameHeader) {
        selectedCol = colIndex;
        selectedRow = -1;
        // 不更新 selectedCellContent，保持編輯中的內容
        return;
      }

      // 只有在點擊不同欄位或不是欄位標題編輯時才結束編輯
      handleEditComplete();
    }

    // 更新選擇狀態和顯示內容
    selectedCol = colIndex;
    selectedRow = -1;
    selectionMode = "column";
    selectedRowRange = new Set();
    selectedColRange = new Set([colIndex]);
    // selectedCellContent 會自動由響應式語句更新
  } // 欄位標題雙擊處理 (進入編輯模式)
  function handleColumnHeaderDblClick(colIndex: number, colName: string) {
    // 設置雙擊標記
    doubleClickInProgress = true;

    // 如果已經在編輯同一個欄位標題，不要重新進入編輯模式
    if (
      editingState.isEditing &&
      editingState.rowIndex === -1 &&
      editingState.colIndex === colIndex
    ) {
      // 清除雙擊標記
      setTimeout(() => {
        doubleClickInProgress = false;
      }, 10);
      return;
    }

    // 如果正在編輯其他元素，先結束編輯
    if (editingState.isEditing) {
      handleEditComplete();
    }

    // 更新選擇狀態，但不更新顯示內容（避免刷新正在輸入的內容）
    selectedCol = colIndex;
    selectedRow = -1;
    // 不更新 selectedCellContent

    // 進入編輯模式
    editingState = {
      tableID,
      rowIndex: -1,
      colIndex,
      colName,
      value: colName,
      isEditing: true,
    };

    // 清除雙擊標記
    setTimeout(() => {
      doubleClickInProgress = false;
    }, 10);
  } // 轉換 nil 值為前端顯示格式
  function formatCellValue(value: any): string {
    // 後端回傳的 nil 值會是 null
    if (value === null || value === undefined) {
      return ".";
    }
    return String(value);
  } // 轉換前端輸入為後端格式
  function parseInputValue(value: string): string {
    // 如果用戶輸入點號，轉換為空字串表示 nil
    if (value === "." || value.trim() === "") {
      return "";
    }
    return value;
  }

  // 編輯完成處理
  async function handleEditComplete() {
    if (!editingState.isEditing) return;

    try {
      if (editingState.rowIndex >= 0) {
        // 更新儲存格值，將用戶輸入的 "." 轉換為空字串
        const processedValue = parseInputValue(editingState.value);
        await UpdateCellValueByID(
          tableID,
          editingState.rowIndex,
          editingState.colIndex,
          processedValue
        );
      } else {
        // 更新欄名
        await UpdateColumnNameByID(
          tableID,
          editingState.colIndex,
          editingState.value
        );
      } // 重新載入資料
      await loadTableData();
    } catch (err) {
      error = `${texts["ui.table.update_failed"] || "更新資料失敗"}: ${err}`;
    } finally {
      // 結束編輯狀態
      editingState = {
        tableID: -1,
        rowIndex: -1,
        colIndex: -1,
        colName: "",
        value: "",
        isEditing: false,
      };
    }
  }

  // 編輯時按下鍵處理
  function handleKeyDown(event: KeyboardEvent) {
    if (event.key === "Enter") {
      event.preventDefault(); // 僅在 Enter 時阻止預設行為
      handleEditComplete();
    } else if (event.key === "Escape") {
      event.preventDefault(); // 僅在 Escape 時阻止預設行為
      // 取消編輯，恢復原值
      editingState = {
        tableID: -1,
        rowIndex: -1,
        colIndex: -1,
        colName: "",
        value: "",
        isEditing: false,
      };
    }
    // 對於其他按鍵，不再呼叫 event.preventDefault()
    // 允許瀏覽器處理正常的文字輸入
  }

  // 將數字索引轉換為字母索引 (A, B, C, ..., AA, AB, ...)
  function indexToLetters(index: number): string {
    if (index < 0) {
      return "A";
    }

    let result = "";
    while (index >= 0) {
      result = String.fromCharCode(65 + (index % 26)) + result;
      index = Math.floor(index / 26) - 1;
      if (index < 0) {
        break;
      }
    }
    return result;
  }
  // 行索引點擊處理 - 選取整行
  function handleRowIndexClick(rowIndex: number) {
    if (editingState.isEditing) {
      handleEditComplete();
    }

    selectionMode = "row";
    selectedRow = rowIndex;
    selectedCol = -1;
    selectedRowRange = new Set([rowIndex]);
    selectedColRange = new Set();
    // selectedCellContent 會自動由響應式語句更新
  }

  // 列索引點擊處理 - 選取整列
  function handleColumnIndexClick(colIndex: number) {
    if (editingState.isEditing) {
      handleEditComplete();
    }

    selectionMode = "column";
    selectedCol = colIndex;
    selectedRow = -1;
    selectedColRange = new Set([colIndex]);
    selectedRowRange = new Set();
    // selectedCellContent 會自動由響應式語句更新
  } // 右鍵菜單處理
  function handleContextMenu(
    event: MouseEvent,
    type: string,
    index?: number,
    rowIndex?: number,
    colIndex?: number
  ) {
    event.preventDefault();

    // 更新菜單標籤（包括貼上狀態）
    updateMenuLabels();

    contextMenuVisible = true;
    contextMenuX = event.clientX;
    contextMenuY = event.clientY;

    // 如果當前是範圍選取模式，使用範圍選取的菜單
    if (
      selectionMode === "range" &&
      rangeSelectStartRow >= 0 &&
      rangeSelectStartCol >= 0
    ) {
      contextMenuType = "range";
    } else {
      contextMenuType = type;
    }

    // 調試信息
    console.log("Mouse position:", {
      clientX: event.clientX,
      clientY: event.clientY,
      pageX: event.pageX,
      pageY: event.pageY,
      screenX: event.screenX,
      screenY: event.screenY,
      offsetX: event.offsetX,
      offsetY: event.offsetY,
    });

    // 設置上下文信息
    contextMenuContext = {
      type: contextMenuType,
      index,
      rowIndex,
      colIndex,
      selectedRow,
      selectedCol,
      tableID,
      rangeSelectStartRow,
      rangeSelectStartCol,
      rangeSelectEndRow,
      rangeSelectEndCol,
    };

    // 根據右鍵類型更新選擇狀態 - 但如果當前是範圍選取模式，則不更新選擇狀態
    if (selectionMode !== "range") {
      if (type === "row" && index !== undefined) {
        handleRowIndexClick(index);
      } else if (type === "column" && index !== undefined) {
        handleColumnIndexClick(index);
      } else if (
        type === "cell" &&
        rowIndex !== undefined &&
        colIndex !== undefined &&
        tableData
      ) {
        const column = tableData.columns[colIndex];
        if (column) {
          const cellValue = tableData.rows[rowIndex]?.cells[column.name];
          const displayValue = formatCellValue(cellValue);
          handleCellClick(rowIndex, colIndex, displayValue);
        }
      }
    }
  }

  // 隱藏右鍵菜單
  function hideContextMenu() {
    contextMenuVisible = false;
  }
  // 點擊文件其他地方時隱藏菜單
  function handleDocumentClick() {
    if (contextMenuVisible) {
      hideContextMenu();
    }
  }

  // 全域鼠標放開事件處理
  function handleGlobalMouseUp() {
    if (isDragging) {
      isDragging = false;
      endRangeSelection();
    }
  } // 右鍵菜單項目處理
  async function handleContextMenuAction(event: CustomEvent) {
    const { action, context } = event.detail;

    console.log("Context menu action:", action, "Context:", context);

    switch (action) {
      case "insertRowAbove":
        console.log(`在第 ${context.rowIndex || context.index} 行上方插入行`);
        break;
      case "insertRowBelow":
        console.log(`在第 ${context.rowIndex || context.index} 行下方插入行`);
        break;
      case "duplicateRow":
        console.log(`複製第 ${context.rowIndex || context.index} 行`);
        // 執行複製整行操作
        if (context.rowIndex !== undefined || context.index !== undefined) {
          const rowIndex = context.rowIndex || context.index;
          // 先設置選取狀態為該行
          selectionMode = "row";
          selectedRow = rowIndex;
          selectedCol = -1;
          selectedRowRange = new Set([rowIndex]);
          selectedColRange = new Set();
          // 然後執行複製
          handleCopy();
        }
        break;
      case "deleteRow":
        console.log(`刪除第 ${context.rowIndex || context.index} 行`);
        break;
      case "insertColumnLeft":
        console.log(`在第 ${context.colIndex || context.index} 欄左邊插入欄`);
        break;
      case "insertColumnRight":
        console.log(`在第 ${context.colIndex || context.index} 欄右邊插入欄`);
        break;
      case "renameColumn":
        console.log(`重命名第 ${context.colIndex || context.index} 欄`);
        // 可以觸發欄位名稱編輯
        break;
      case "duplicateColumn":
        console.log(`複製第 ${context.colIndex || context.index} 欄`);
        // 執行複製整欄操作
        if (context.colIndex !== undefined || context.index !== undefined) {
          const colIndex = context.colIndex || context.index;
          // 先設置選取狀態為該欄
          selectionMode = "column";
          selectedCol = colIndex;
          selectedRow = -1;
          selectedColRange = new Set([colIndex]);
          selectedRowRange = new Set();
          // 然後執行複製
          handleCopy();
        }
        break;
      case "deleteColumn":
        console.log(`刪除第 ${context.colIndex || context.index} 欄`);
        break;
      case "copy":
        console.log(`複製儲存格 (${context.rowIndex}, ${context.colIndex})`);
        handleCopy();
        break;
      case "paste":
        // 根據上下文類型決定貼上方式
        if (
          context.type === "row" &&
          (context.rowIndex !== undefined || context.index !== undefined)
        ) {
          const rowIndex = context.rowIndex || context.index;
          console.log(`貼上到第 ${rowIndex} 列`);
          await handlePasteToRow(rowIndex);
        } else if (
          context.type === "column" &&
          (context.colIndex !== undefined || context.index !== undefined)
        ) {
          const colIndex = context.colIndex || context.index;
          console.log(`貼上到第 ${colIndex} 欄`);
          await handlePasteToColumn(colIndex);
        } else {
          // 預設為儲存格貼上
          console.log(
            `貼上到儲存格 (${context.rowIndex}, ${context.colIndex})`
          );
          await handlePaste();
        }
        break;
      case "clear":
        console.log(`清除儲存格 (${context.rowIndex}, ${context.colIndex})`);
        // 實現清除功能
        if (context.rowIndex !== undefined && context.colIndex !== undefined) {
          await UpdateCellValueByID(
            tableID,
            context.rowIndex,
            context.colIndex,
            ""
          );
          // 重新載入資料
          await loadTableData();
        }
        break;
      default:
        console.log("未知的菜單動作:", action);
    }
    hideContextMenu();
  }
</script>

<div class="data-table-container">
  {#if loading}
    <div class="loading">{texts["ui.table.loading"] || "載入中..."}</div>
  {:else if error}
    <div class="error">{error}</div>
  {:else if tableData}
    <div class="table-wrapper" style="--table-scale: {tableScale};">
      <table class="data-table">
        <thead>
          <!-- 欄位索引行 (A, B, C, ...) -->
          <tr>
            <!-- 空白頂角儲存格 -->
            <th class="corner-cell corner-index"></th>
            <!-- 欄位索引 -->
            {#each tableData.columns as column, colIndex}
              <th
                class="column-index"
                class:selected={colIndex === selectedCol ||
                  (selectionMode === "column" &&
                    selectedColRange.has(colIndex)) ||
                  (selectionMode === "cell" && colIndex === selectedCol)}
                on:click={() => handleColumnIndexClick(colIndex)}
                on:contextmenu={(e) =>
                  handleContextMenu(e, "column", colIndex)}
              >
                {indexToLetters(colIndex)}
              </th>
            {/each}
          </tr>

          <!-- 欄位名稱行 -->
          <tr>
            <!-- 空白頂角儲存格 -->
            <th class="corner-cell corner-header"></th>

            <!-- 欄位標題 -->
            {#each tableData.columns as column, colIndex}
              <th
                class="column-header"
                class:selected={colIndex === selectedCol ||
                  (selectionMode === "column" &&
                    selectedColRange.has(colIndex)) ||
                  (selectionMode === "cell" && colIndex === selectedCol)}
                on:click={() =>
                  handleColumnHeaderClick(colIndex, column.name)}
                on:dblclick={() =>
                  handleColumnHeaderDblClick(colIndex, column.name)}
                on:contextmenu={(e) =>
                  handleContextMenu(e, "column", colIndex)}
              >
                {#if editingState.isEditing && editingState.rowIndex === -1 && editingState.colIndex === colIndex}
                  <input
                    type="text"
                    bind:value={editingState.value}
                    on:keydown={handleKeyDown}
                    on:blur={handleEditComplete}
                    class="editor"
                    bind:this={editInput}
                  />
                {:else}
                  {column.name}
                {/if}
              </th>
            {/each}
          </tr>
        </thead>
        <tbody>
          {#each tableData.rows as row, rowIndex}
            <tr
              class:selected-row={rowIndex === selectedRow ||
                (selectionMode === "row" && selectedRowRange.has(rowIndex))}
            >
              <!-- 行標識 -->
              <td
                class="row-header"
                class:selected={rowIndex === selectedRow ||
                  (selectionMode === "row" && selectedRowRange.has(rowIndex))}
                on:click={() => handleRowIndexClick(rowIndex)}
                on:contextmenu={(e) => handleContextMenu(e, "row", rowIndex)}
              >
                {rowIndex + 1}
              </td>
              <!-- 儲存格資料 -->
              {#each tableData.columns as column, colIndex}
                {@const cellValue = row.cells[column.name]}
                {@const displayValue = formatCellValue(cellValue)}
                {@const isInRange =
                  selectionMode === "range" &&
                  rangeSelectStartRow >= 0 &&
                  rangeSelectStartCol >= 0 &&
                  rangeSelectEndRow >= 0 &&
                  rangeSelectEndCol >= 0 &&
                  rowIndex >=
                    Math.min(rangeSelectStartRow, rangeSelectEndRow) &&
                  rowIndex <=
                    Math.max(rangeSelectStartRow, rangeSelectEndRow) &&
                  colIndex >=
                    Math.min(rangeSelectStartCol, rangeSelectEndCol) &&
                  colIndex <=
                    Math.max(rangeSelectStartCol, rangeSelectEndCol)}
                <td
                  class="cell"
                  class:selected-cell={rowIndex === selectedRow &&
                    colIndex === selectedCol}
                  class:selected-col={colIndex === selectedCol ||
                    (selectionMode === "column" &&
                      selectedColRange.has(colIndex)) ||
                    (selectionMode === "cell" && colIndex === selectedCol)}
                  class:selected-row-cell={rowIndex === selectedRow ||
                    (selectionMode === "row" &&
                      selectedRowRange.has(rowIndex))}
                  class:range-selected={isInRange}
                  class:nil-value={cellValue === null ||
                    cellValue === undefined}
                  on:click={(e) =>
                    handleCellClick(rowIndex, colIndex, displayValue, e)}
                  on:mousedown={(e) =>
                    handleCellMouseDown(rowIndex, colIndex, e)}
                  on:mouseup={() => handleCellMouseUp()}
                  on:dblclick={() =>
                    handleCellDblClick(
                      rowIndex,
                      colIndex,
                      column.name,
                      cellValue
                    )}
                  on:contextmenu={(e) =>
                    handleContextMenu(
                      e,
                      "cell",
                      undefined,
                      rowIndex,
                      colIndex
                    )}
                  on:mouseenter={() =>
                    handleCellMouseEnter(rowIndex, colIndex)}
                  on:keydown={(e) => {
                    if (e.key === "Enter" || e.key === " ") {
                      handleCellClick(rowIndex, colIndex, displayValue, e);
                    }
                  }}
                  tabindex="0"
                  role="gridcell"
                  >{#if editingState.isEditing && editingState.rowIndex === rowIndex && editingState.colIndex === colIndex}
                    <input
                      type="text"
                      bind:value={editingState.value}
                      on:keydown={handleKeyDown}
                      on:blur={handleEditComplete}
                      class="editor"
                      bind:this={editInput}
                    />
                  {:else}
                    {displayValue}
                  {/if}
                </td>
              {/each}
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
    <!-- /.table-wrapper -->
    <div class="table-controls">
      <div class="selected-content">
        <strong>{texts["ui.table.selected_content"] || "選中內容"}:</strong>
        {selectedCellContent}
      </div>
      <div class="zoom-control">
        <label for="table-zoom">{texts["ui.table.zoom"] || "縮放"}:</label>
        <input
          type="range"
          id="table-zoom"
          min="0.5"
          max="2"
          step="0.1"
          bind:value={tableScale}
        />
        <span>{Math.round(tableScale * 100)}%</span>
      </div>
    </div>
  {:else}
    <div class="no-data">{texts["ui.table.no_data"] || "無資料可顯示"}</div>
  {/if}
</div>

<!-- 右鍵菜單組件 -->
<ContextMenu
  visible={contextMenuVisible}
  x={contextMenuX}
  y={contextMenuY}
  type={contextMenuType}
  menuConfig={contextMenuConfig}
  context={contextMenuContext}
  on:action={handleContextMenuAction}
  on:close={hideContextMenu}
/>

<style>
  .data-table-container {
    width: 100%;
    height: 100%;
    overflow: hidden; /* 改為 hidden，讓 table-wrapper 處理滾動 */
    display: flex;
    flex-direction: column;
    font-family:
      "Nunito",
      -apple-system,
      BlinkMacSystemFont,
      "Segoe UI",
      Roboto,
      Oxygen,
      Ubuntu,
      Cantarell,
      "Open Sans",
      "Helvetica Neue",
      sans-serif;
    background: rgba(255, 255, 255, 0.95);
    border-radius: var(--radius-large);
    box-shadow: var(--shadow-2);
    backdrop-filter: blur(20px);
    border: 1px solid rgba(255, 255, 255, 0.2);
  }

  .loading,
  .error,
  .no-data {
    padding: var(--spacing-xl);
    text-align: center;
    display: flex;
    align-items: center;
    justify-content: center;
    flex-direction: column;
    height: 100%;
    background: linear-gradient(
      135deg,
      rgba(255, 255, 255, 0.9),
      rgba(248, 250, 252, 0.9)
    );
    border-radius: var(--radius-large);
  }

  .loading {
    color: var(--primary-color);
    font-size: 1.1rem;
    font-weight: 500;
  }

  .loading::before {
    content: "";
    width: 40px;
    height: 40px;
    border: 3px solid var(--primary-light);
    border-top: 3px solid var(--primary-color);
    border-radius: 50%;
    animation: spin 1s linear infinite;
    margin-bottom: var(--spacing-md);
  }

  @keyframes spin {
    0% {
      transform: rotate(0deg);
    }
    100% {
      transform: rotate(360deg);
    }
  }

  .error {
    color: var(--error-color);
    font-weight: 500;
  }
  .no-data {
    color: var(--text-secondary);
    font-style: italic;
  }
  .table-container-scaled {
    /* 縮放容器，確保縮放不會影響內部的 sticky positioning */
    transform-origin: top left;
    /* 動態調整寬度和高度以適應縮放 */
    width: calc(100% / var(--table-scale, 1));
    height: calc(100% / var(--table-scale, 1));
    overflow: hidden;
    flex: 1; /* 填充剩餘垂直空間，控制列會自動占用其需要的空間 */
    /* 確保縮放後仍能正確顯示滾動條 */
    min-width: 100%;
    /* 添加過渡動畫以平滑縮放 */
    transition: transform 0.2s ease-out;
  }
  .table-wrapper {
    /* 使用 CSS 變數來控制尺寸 */
    overflow: auto; /* 確保 wrapper 可以滾動 */
    margin: var(--spacing-sm);
    border-radius: var(--radius-medium);
    box-shadow: inset 0 2px 4px rgba(0, 0, 0, 0.06);
    background: rgba(255, 255, 255, 0.8);
    backdrop-filter: blur(10px);
    position: relative;
    /* 為縮放調整容器尺寸 - 修改為自適應高度 */
    width: calc(100% - 2 * var(--spacing-sm));
    flex: 1; /* 讓表格區域填滿剩餘空間 */
    /* 確保在縮放時能正確顯示內容 */
    box-sizing: border-box;
    /* 設定最小高度以避免被壓縮 */
    min-height: 300px;
    /* 使用 CSS 變數控制縮放 */
    font-size: calc(0.9rem * var(--table-scale, 1));
  }

  .data-table {
    border-collapse: separate;
    border-spacing: 0;
    table-layout: auto; /* 改為 auto 讓儲存格保持固定寬度 */
    background: transparent;
    min-width: max-content; /* 確保表格內容不會被過度壓縮 */
    user-select: none; /* 防止文字選取 */
    /* 移除 width: 100% 讓表格寬度由內容決定 */
    width: auto;
  }

  th,
  td {
    padding: calc(var(--spacing-sm) * var(--table-scale, 1)) calc(var(--spacing-md) * var(--table-scale, 1));
    text-align: left;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    width: calc(140px * var(--table-scale, 1));
    min-width: calc(140px * var(--table-scale, 1));
    max-width: calc(140px * var(--table-scale, 1));
    border: none;
    border-bottom: 1px solid rgba(0, 0, 0, 0.06);
    border-right: 1px solid rgba(0, 0, 0, 0.06);
    position: relative;
    transition: all var(--transition-fast);
    font-size: inherit;
  }

  .corner-cell {
    background: linear-gradient(135deg, #f8fafc, #e2e8f0);
    width: calc(60px * var(--table-scale, 1));
    min-width: calc(60px * var(--table-scale, 1));
    max-width: calc(60px * var(--table-scale, 1));
    position: sticky;
    left: 0;
    z-index: 15;
    box-shadow: 2px 0 4px rgba(0, 0, 0, 0.1);
    border-right: 2px solid rgba(0, 0, 0, 0.1) !important;
  }

  .corner-index {
    top: 0;
    border-radius: var(--radius-medium) 0 0 0;
  }

  .corner-header {
    top: calc(40px * var(--table-scale, 1));
  }

  .column-index {
    background: linear-gradient(180deg, #f1f5f9, #e2e8f0);
    position: sticky;
    top: 0;
    z-index: 12;
    font-weight: 600;
    text-align: center;
    color: var(--text-secondary);
    font-size: calc(0.85rem * var(--table-scale, 1));
    letter-spacing: 0.5px;
    height: calc(40px * var(--table-scale, 1));
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.08);
    border-bottom: 2px solid rgba(0, 0, 0, 0.1) !important;
  }

  .column-index.selected {
    background: linear-gradient(
      180deg,
      var(--primary-light),
      var(--primary-color)
    );
    color: var(--text-on-primary);
    box-shadow: 0 4px 8px rgba(25, 118, 210, 0.3);
    transform: translateY(-1px);
  }

  .column-header {
    background: linear-gradient(
      180deg,
      rgba(225, 238, 255, 0.9),
      rgba(191, 219, 254, 0.8)
    );
    position: sticky;
    top: calc(40px * var(--table-scale, 1));
    z-index: 11;
    font-weight: 600;
    height: calc(40px * var(--table-scale, 1));
    color: var(--text-primary);
    backdrop-filter: blur(10px);
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.08);
    border-bottom: 2px solid rgba(0, 0, 0, 0.1) !important;
  }

  .column-header.selected {
    background: linear-gradient(
      180deg,
      var(--secondary-light),
      var(--secondary-color)
    );
    color: var(--text-primary);
    box-shadow: 0 4px 8px rgba(3, 218, 198, 0.3);
    transform: translateY(-1px);
  }

  .row-header {
    background: linear-gradient(90deg, #f8fafc, #e2e8f0);
    position: sticky;
    left: 0;
    z-index: 9;
    width: calc(60px * var(--table-scale, 1));
    min-width: calc(60px * var(--table-scale, 1));
    max-width: calc(60px * var(--table-scale, 1));
    font-weight: 600;
    text-align: center;
    color: var(--text-secondary);
    font-size: calc(0.9rem * var(--table-scale, 1));
    box-shadow: 2px 0 4px rgba(0, 0, 0, 0.1);
    border-right: 2px solid rgba(0, 0, 0, 0.1) !important;
  }

  .row-header.selected {
    background: linear-gradient(
      90deg,
      var(--primary-light),
      var(--primary-color)
    );
    color: var(--text-on-primary);
    box-shadow: 4px 0 8px rgba(25, 118, 210, 0.3);
    transform: translateX(1px);
  }

  .cell {
    position: relative;
    background: rgba(255, 255, 255, 0.7);
    cursor: pointer;
    font-size: inherit;
  }
  .cell:hover:not(.selected-col):not(.selected-row-cell):not(.selected-cell) {
    background: rgba(25, 118, 210, 0.08);
    transform: scale(1.02);
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  }

  .selected-row {
    background: linear-gradient(
      90deg,
      rgba(25, 118, 210, 0.12),
      rgba(25, 118, 210, 0.06)
    ) !important;
  }
  .selected-col {
    background: linear-gradient(
      180deg,
      rgba(3, 218, 198, 0.35),
      rgba(3, 218, 198, 0.2)
    ) !important;
    box-shadow: inset 4px 0 0 rgba(3, 218, 198, 0.8) !important;
    position: relative;
    z-index: 3;
  }

  .selected-row-cell {
    background: linear-gradient(
      90deg,
      rgba(25, 118, 210, 0.12),
      rgba(25, 118, 210, 0.06)
    ) !important;
  }

  .selected-cell {
    background: linear-gradient(
      135deg,
      rgba(25, 118, 210, 0.15),
      rgba(3, 218, 198, 0.15)
    ) !important;
    box-shadow: inset 0 0 0 2px var(--primary-color) !important;
    position: relative;
    z-index: 5;
  }

  .range-selected {
    background: linear-gradient(
      135deg,
      rgba(25, 118, 210, 0.2),
      rgba(3, 218, 198, 0.1)
    ) !important;
    box-shadow: inset 0 0 0 2px rgba(25, 118, 210, 0.6) !important;
    position: relative;
    z-index: 4;
  }

  /* 當整欄被選中且在範圍選取內時，保持欄位高亮但加上範圍選取的邊框 */
  .selected-col.range-selected {
    background: linear-gradient(
      180deg,
      rgba(3, 218, 198, 0.45),
      rgba(3, 218, 198, 0.3)
    ) !important;
    box-shadow:
      inset 4px 0 0 rgba(3, 218, 198, 0.8),
      inset 0 0 0 2px rgba(25, 118, 210, 0.6) !important;
    position: relative;
    z-index: 6;
  }

  /* 防止在拖拽時選取文字 */
  .cell.range-selected {
    cursor: cell;
  }
  .selected-row .cell:not(.selected-col) {
    background: linear-gradient(
      90deg,
      rgba(25, 118, 210, 0.08),
      rgba(25, 118, 210, 0.04)
    ) !important;
  }

  .selected-row .selected-col {
    background: linear-gradient(
      135deg,
      rgba(25, 118, 210, 0.15),
      rgba(3, 218, 198, 0.15)
    ) !important;
    box-shadow:
      0 0 0 2px var(--primary-color),
      0 4px 12px rgba(25, 118, 210, 0.3) !important;
    border-radius: var(--radius-small) !important;
  }

  .nil-value {
    color: var(--text-hint);
    font-style: italic;
    background: linear-gradient(
      135deg,
      rgba(0, 0, 0, 0.02),
      rgba(0, 0, 0, 0.01)
    ) !important;
  }

  .nil-value::before {
    content: "∅";
    opacity: 0.3;
    margin-right: var(--spacing-xs);
  }

  .editor {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    padding: calc(var(--spacing-sm) * var(--table-scale, 1));
    border: 2px solid var(--primary-color);
    border-radius: var(--radius-small);
    box-sizing: border-box;
    outline: none;
    background: rgba(255, 255, 255, 0.95);
    backdrop-filter: blur(10px);
    font-family: inherit;
    font-size: inherit;
    color: var(--text-primary);
    box-shadow: 0 4px 12px rgba(25, 118, 210, 0.3);
    z-index: 10;
  }

  .editor:focus {
    box-shadow:
      0 0 0 3px rgba(25, 118, 210, 0.2),
      0 4px 12px rgba(25, 118, 210, 0.4);
  }
  .table-controls {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: var(--spacing-sm) var(--spacing-md);
    border-top: 1px solid rgba(0, 0, 0, 0.06);
    background: rgba(248, 250, 252, 0.9);
    position: relative; /* 新增，為 z-index 生效 */
    z-index: 20; /* 新增，確保在縮放內容之上 */
    height: 60px; /* 固定高度，確保一致性 */
    flex-shrink: 0; /* 防止被壓縮 */
    min-height: 60px; /* 最小高度保證 */
  }

  .selected-content {
    flex-grow: 1;
    font-size: 0.9rem;
    color: var(--text-secondary);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    margin-right: var(--spacing-md);
  }

  .zoom-control {
    display: flex;
    align-items: center;
    gap: var(--spacing-xs);
  }

  .zoom-control label {
    font-size: 0.9rem;
    color: var(--text-secondary);
  }

  .zoom-control input[type="range"] {
    width: 200px; /* 增加寬度以便更精確控制 */
    -webkit-appearance: none; /* 移除預設樣式 (Chrome, Safari, Opera) */
    appearance: none;
    height: 10px; /* 軌道高度 */
    background: #e0e5ec; /* 新擬物背景色 */
    border-radius: 5px;
    outline: none;
    box-shadow:
      inset 3px 3px 6px #b8bec7,
      /* 內陰影 - 暗 */ inset -3px -3px 6px #ffffff; /* 內陰影 - 亮 */
    transition: box-shadow 0.15s ease-in-out; /* 添加過渡效果 */
  }

  /* 確保軌道在 active 狀態下樣式不變 */
  .zoom-control input[type="range"]:active {
    background: #e0e5ec; /* 保持背景不變 */
    box-shadow:
      inset 3px 3px 6px #b8bec7,
      inset -3px -3px 6px #ffffff; /* 保持陰影不變 */
  }

  /* Webkit (Chrome, Safari, Opera) 瀏覽器的滑塊樣式 */
  .zoom-control input[type="range"]::-webkit-slider-thumb {
    -webkit-appearance: none;
    appearance: none;
    width: 20px; /* 滑塊寬度 */
    height: 20px; /* 滑塊高度 */
    background: #e0e5ec; /* 滑塊背景色 */
    border-radius: 50%; /* 圓形滑塊 */
    cursor: pointer;
    border: 1px solid #c8cdd3;
    box-shadow:
      3px 3px 6px #b8bec7,
      /* 外陰影 - 暗 */ -3px -3px 6px #ffffff; /* 外陰影 - 亮 */
    transition: background-color 0.15s ease-in-out;
  }

  .zoom-control input[type="range"]::-webkit-slider-thumb:active {
    background-color: #d1d9e6; /* 輕微改變背景色以示選中，移除內陰影 */
  }

  /* Mozilla Firefox 瀏覽器的滑塊樣式 */
  .zoom-control input[type="range"]::-moz-range-thumb {
    width: 18px; /* 滑塊寬度 */
    height: 18px; /* 滑塊高度 */
    background: #e0e5ec;
    border-radius: 50%;
    cursor: pointer;
    border: 1px solid #c8cdd3;
    box-shadow:
      3px 3px 6px #b8bec7,
      -3px -3px 6px #ffffff;
    transition: background-color 0.15s ease-in-out;
  }

  .zoom-control input[type="range"]::-moz-range-thumb:active {
    /* Add a dummy property to avoid empty ruleset */
    border: none;
  }

  /* Mozilla Firefox 瀏覽器的軌道樣式 (可選，如果需要更細緻的控制) */
  .zoom-control input[type="range"]::-moz-range-track {
    width: 100%;
    height: 10px;
    background: #e0e5ec;
    border-radius: 5px;
    box-shadow:
      inset 3px 3px 6px #b8bec7,
      inset -3px -3px 6px #ffffff;
    border: none; /* 確保無邊框影響 */
  }

  /* 確保 Firefox 軌道在 active 狀態下樣式不變 */
  .zoom-control input[type="range"]:active::-moz-range-track {
    background: #e0e5ec; /* 保持背景不變 */
    box-shadow:
      inset 3px 3px 6px #b8bec7,
      inset -3px -3px 6px #ffffff; /* 保持陰影不變 */
  }

  .zoom-control span {
    font-size: 0.9rem;
    color: var(--text-primary);
    /* min-width: 35px; */ /* 移除 min-width */
    width: 45px; /* 設定固定寬度以容納三位數百分比 */
    text-align: right;
    display: inline-block; /* 確保寬度生效 */
  } /* 滾動條樣式 - 根據縮放比例調整 */
  .table-wrapper::-webkit-scrollbar-track {
    background: rgba(0, 0, 0, 0.05);
    border-radius: var(--radius-medium);
  }

  .table-wrapper::-webkit-scrollbar-thumb {
    background: linear-gradient(
      135deg,
      var(--primary-color),
      var(--primary-light)
    );
    border-radius: var(--radius-medium);
    border: 2px solid transparent;
    background-clip: content-box;
  }

  .table-wrapper::-webkit-scrollbar-thumb:hover {
    background: linear-gradient(
      135deg,
      var(--primary-dark),
      var(--primary-color)
    );
    background-clip: content-box;
  }

  .table-wrapper::-webkit-scrollbar-corner {
    background: rgba(0, 0, 0, 0.05);
  }

  /* 響應式設計 */
  @media (max-width: 768px) {
    th,
    td {
      width: 120px;
      min-width: 120px;
      max-width: 120px;
      padding: var(--spacing-xs) var(--spacing-sm);
      font-size: 0.8rem;
    }

    .corner-cell,
    .row-header {
      width: 50px;
      min-width: 50px;
      max-width: 50px;
    }

    .column-index,
    .column-header {
      height: 35px;
    }
  }
</style>
