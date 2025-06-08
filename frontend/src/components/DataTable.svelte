<script lang="ts">
  import { onMount, createEventDispatcher } from "svelte";
  import { TableData, EditingStateByID } from "../types/datatable";
  import {
    GetTableDataByID,
    UpdateCellValueByID,
    UpdateColumnNameByID,
  } from "../../wailsjs/go/main/App";
  import ContextMenu from "./ContextMenu.svelte";
  import type { ContextMenuConfig } from "../types/contextMenu";

  // 組件屬性
  export let tableID: number;
  export let tableKey: number = 0; // 用於強制重新載入的 key

  // 創建事件分發器
  const dispatch = createEventDispatcher();
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
  };
  // 選中狀態
  let selectedRow = -1;
  let selectedCol = -1;
  let selectedCellContent = "";

  // 選擇模式：'cell' | 'row' | 'column'
  let selectionMode = "cell";
  let selectedRowRange = new Set(); // 選中的行範圍
  let selectedColRange = new Set(); // 選中的列範圍

  // 右鍵菜單狀態
  let contextMenuVisible = false;
  let contextMenuX = 0;
  let contextMenuY = 0;
  let contextMenuType = ""; // 'row' | 'column' | 'cell'
  let contextMenuContext = {}; // 上下文信息

  // 右鍵菜單配置
  const contextMenuConfig: ContextMenuConfig = {
    row: [
      { id: "insertRowAbove", label: "在上方插入行", icon: "⬆️" },
      { id: "insertRowBelow", label: "在下方插入行", icon: "⬇️" },
      { id: "separator1", type: "separator" },
      { id: "duplicateRow", label: "複製行", icon: "📋" },
      { id: "separator2", type: "separator" },
      { id: "deleteRow", label: "刪除行", icon: "🗑️", danger: true },
    ],
    column: [
      { id: "insertColumnLeft", label: "在左邊插入欄", icon: "⬅️" },
      { id: "insertColumnRight", label: "在右邊插入欄", icon: "➡️" },
      { id: "separator1", type: "separator" },
      { id: "renameColumn", label: "重新命名欄位", icon: "✏️" },
      { id: "duplicateColumn", label: "複製欄位", icon: "📋" },
      { id: "separator2", type: "separator" },
      { id: "deleteColumn", label: "刪除欄位", icon: "🗑️", danger: true },
    ],
    cell: [
      { id: "copy", label: "複製", icon: "📋" },
      { id: "paste", label: "貼上", icon: "📄", disabled: true }, // 可以根據剪貼簿狀態動態設置
      { id: "separator1", type: "separator" },
      { id: "clear", label: "清除內容", icon: "🧹" },
      { id: "separator2", type: "separator" },
      { id: "insertRowAbove", label: "在上方插入行", icon: "⬆️" },
      { id: "insertRowBelow", label: "在下方插入行", icon: "⬇️" },
      { id: "insertColumnLeft", label: "在左邊插入欄", icon: "⬅️" },
      { id: "insertColumnRight", label: "在右邊插入欄", icon: "➡️" },
    ],
  };

  // 防止雙擊時觸發點擊的標記
  let doubleClickInProgress = false;
  // 編輯輸入元素引用
  let editInput: HTMLInputElement; // 當進入編輯模式時，設置焦點
  $: if (editingState.isEditing && editInput) {
    setTimeout(() => {
      editInput.focus();
      editInput.select();
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

  // 更新選中內容顯示的函數
  function updateSelectedCellContent() {
    if (!tableData) return;

    if (selectionMode === "row" && selectedRow >= 0) {
      selectedCellContent = `第 ${selectedRow + 1} 行`;
    } else if (selectionMode === "column" && selectedCol >= 0) {
      selectedCellContent = `${indexToLetters(selectedCol)} 欄`;
    } else if (
      selectionMode === "cell" &&
      selectedRow >= 0 &&
      selectedCol >= 0
    ) {
      const column = tableData.columns[selectedCol];
      if (column) {
        const cellValue = tableData.rows[selectedRow]?.cells[column.name];
        const displayValue = formatCellValue(cellValue);
        selectedCellContent = `${indexToLetters(selectedCol)}${selectedRow + 1}: ${displayValue}`;
      }
    } else {
      selectedCellContent = "";
    }
  }
  onMount(async () => {
    lastTableID = tableID;
    lastTableKey = tableKey;
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
  }
  // 載入表格資料
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
      總行數: rowCount.toString(),
      總欄數: colCount.toString(),
      總儲存格: totalCells.toString(),
      數值欄數: numericCols.toString(),
    };
  } // 儲存格點擊處理
  function handleCellClick(rowIndex: number, colIndex: number, value: string) {
    // 如果正在雙擊過程中，忽略點擊事件
    if (doubleClickInProgress) {
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

    // 更新選擇狀態為儲存格模式
    selectionMode = "cell";
    selectedRow = rowIndex;
    selectedCol = colIndex;
    selectedRowRange = new Set();
    selectedColRange = new Set();
    // selectedCellContent 會自動由響應式語句更新
  } // 儲存格雙擊處理 (進入編輯模式)
  function handleCellDblClick(
    rowIndex: number,
    colIndex: number,
    colName: string,
    value: string
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

    // 進入編輯模式
    editingState = {
      tableID,
      rowIndex,
      colIndex,
      colName,
      value,
      isEditing: true,
    };

    // 清除雙擊標記
    setTimeout(() => {
      doubleClickInProgress = false;
    }, 10);
  } // 欄位標題點擊處理
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
  }
  // 轉換 nil 值為前端顯示格式
  function formatCellValue(value: any): string {
    if (value === null || value === undefined) {
      return ".";
    }
    return String(value);
  }
  // 轉換前端輸入為後端格式
  function parseInputValue(value: string): string {
    // 直接傳送用戶輸入的值，讓後端決定如何處理
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
      }

      // 重新載入資料
      await loadTableData();
    } catch (err) {
      error = `更新資料失敗: ${err}`;
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

  // 編輯內容變更處理
  function handleEditChange(event: Event) {
    const target = event.target as HTMLInputElement;
    editingState.value = target.value;
  }
  // 編輯時按下 Enter 鍵處理
  function handleKeyDown(event: KeyboardEvent) {
    if (event.key === "Enter") {
      handleEditComplete();
    }
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
  }
  // 右鍵菜單處理
  function handleContextMenu(
    event: MouseEvent,
    type: string,
    index?: number,
    rowIndex?: number,
    colIndex?: number
  ) {
    event.preventDefault();
    contextMenuVisible = true;
    contextMenuX = event.clientX;
    contextMenuY = event.clientY;
    contextMenuType = type;

    // 設置上下文信息
    contextMenuContext = {
      type,
      index,
      rowIndex,
      colIndex,
      selectedRow,
      selectedCol,
      tableID,
    }; // 根據右鍵類型更新選擇狀態
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

  // 隱藏右鍵菜單
  function hideContextMenu() {
    contextMenuVisible = false;
  }

  // 點擊文件其他地方時隱藏菜單
  function handleDocumentClick() {
    if (contextMenuVisible) {
      hideContextMenu();
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
        break;
      case "deleteColumn":
        console.log(`刪除第 ${context.colIndex || context.index} 欄`);
        break;
      case "copy":
        console.log(`複製儲存格 (${context.rowIndex}, ${context.colIndex})`);
        // 實現複製功能
        break;
      case "paste":
        console.log(`貼上到儲存格 (${context.rowIndex}, ${context.colIndex})`);
        // 實現貼上功能
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
    <div class="loading">載入中...</div>
  {:else if error}
    <div class="error">{error}</div>
  {:else if tableData}
    <div class="table-wrapper">
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
                    selectedColRange.has(colIndex))}
                on:click={() => handleColumnIndexClick(colIndex)}
                on:contextmenu={(e) => handleContextMenu(e, "column", colIndex)}
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
                class:selected={colIndex === selectedCol}
                on:click={() => handleColumnHeaderClick(colIndex, column.name)}
                on:dblclick={() =>
                  handleColumnHeaderDblClick(colIndex, column.name)}
              >
                {#if editingState.isEditing && editingState.rowIndex === -1 && editingState.colIndex === colIndex}
                  <input
                    type="text"
                    bind:value={editingState.value}
                    on:change={handleEditChange}
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
                <td
                  class="cell"
                  class:selected-cell={rowIndex === selectedRow &&
                    colIndex === selectedCol}
                  class:selected-col={colIndex === selectedCol ||
                    (selectionMode === "column" &&
                      selectedColRange.has(colIndex))}
                  class:selected-row-cell={rowIndex === selectedRow ||
                    (selectionMode === "row" && selectedRowRange.has(rowIndex))}
                  class:nil-value={cellValue === null ||
                    cellValue === undefined}
                  on:click={() =>
                    handleCellClick(rowIndex, colIndex, displayValue)}
                  on:dblclick={() =>
                    handleCellDblClick(
                      rowIndex,
                      colIndex,
                      column.name,
                      displayValue
                    )}
                  on:contextmenu={(e) =>
                    handleContextMenu(e, "cell", undefined, rowIndex, colIndex)}
                  on:keydown={(e) => {
                    if (e.key === "Enter" || e.key === " ") {
                      handleCellClick(rowIndex, colIndex, displayValue);
                    }
                  }}
                  tabindex="0"
                  role="gridcell"
                  >{#if editingState.isEditing && editingState.rowIndex === rowIndex && editingState.colIndex === colIndex}
                    <input
                      type="text"
                      bind:value={editingState.value}
                      on:change={handleEditChange}
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
    {#if selectedCellContent}
      <div class="selected-content">
        <strong>選中內容:</strong>
        {selectedCellContent}
      </div>
    {/if}
  {:else}
    <div class="no-data">無資料可顯示</div>
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
    overflow: hidden;
    display: flex;
    flex-direction: column;
    font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Oxygen,
      Ubuntu, Cantarell, "Open Sans", "Helvetica Neue", sans-serif;
  }

  .loading,
  .error,
  .no-data {
    padding: 20px;
    text-align: center;
  }

  .error {
    color: #d32f2f;
  }

  .table-wrapper {
    flex: 1;
    overflow: auto;
    border: 1px solid #ddd;
    border-radius: 4px;
  }
  .data-table {
    border-collapse: collapse;
    table-layout: fixed;
    /* 移除 min-width，讓表格可以超出容器 */
  }

  th,
  td {
    padding: 8px 12px;
    border: 1px solid #ddd;
    text-align: left;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    width: 120px; /* 固定寬度 */
    min-width: 120px;
    max-width: 120px;
  }
  .corner-cell {
    background-color: #f0f0f0;
    width: 50px;
    min-width: 50px;
    max-width: 50px;
    position: sticky;
    left: 0;
    z-index: 15; /* 最高層級，確保在其他sticky元素之上 */
  }

  .corner-index {
    top: 0;
  }

  .corner-header {
    top: 35px; /* 與索引行高度一致 */
  }

  .column-index {
    background-color: #f8f9fa;
    position: sticky;
    top: 0;
    z-index: 12;
    font-weight: 500;
    text-align: center;
    color: #666;
    font-size: 0.9rem;
    transition: background-color 0.2s;
    height: 35px; /* 設定固定高度 */
    padding: 6px 12px; /* 調整padding */
  }

  .column-index.selected {
    background-color: #c8e6c9;
  }

  .column-header {
    background-color: #e1eeff;
    position: sticky;
    top: 35px; /* 與索引行高度一致 */
    z-index: 11;
    font-weight: 600;
    transition: background-color 0.2s;
    height: 35px; /* 設定固定高度 */
    padding: 6px 12px; /* 調整padding */
  }

  .column-header.selected {
    background-color: #c8e6c9;
  }
  .row-header {
    background-color: #f0f0f0;
    position: sticky;
    left: 0;
    z-index: 9;
    width: 50px;
    min-width: 50px;
    max-width: 50px;
    transition: background-color 0.2s;
  }

  .row-header.selected {
    background-color: #c8e6c9;
  }
  .selected-row {
    background-color: rgba(200, 230, 201, 0.3) !important;
  }

  .selected-col {
    background-color: rgba(200, 230, 201, 0.3) !important;
  }

  .selected-row-cell {
    background-color: rgba(200, 230, 201, 0.3) !important;
  }
  .selected-cell {
    background-color: rgba(123, 31, 162, 0.15) !important; /* 更淡的紫色背景 */
    border: 2px solid rgb(94, 23, 125) !important; /* 更淡的紫色邊框 */
    box-sizing: border-box;
  }

  /* 確保選中的行中所有儲存格都高亮 */
  .selected-row .cell {
    background-color: rgba(200, 230, 201, 0.3) !important;
  }

  /* 確保選中的欄中所有儲存格都高亮 */
  .selected-col {
    background-color: rgba(200, 230, 201, 0.3) !important;
  }

  /* 被選中的儲存格具有最高優先級 - 更淡的紫色背景 */
  .selected-row .selected-col {
    background-color: rgba(123, 31, 162, 0.15) !important; /* 更淡的紫色背景 */
    border: 2px solid rgba(123, 31, 162, 0.5) !important; /* 更淡的紫色邊框 */
    box-sizing: border-box;
  }

  .cell {
    position: relative;
    transition: background-color 0.2s;
  }
  .cell:hover {
    background-color: #f9f9f9;
  }

  .nil-value {
    color: #999;
    font-style: italic;
    background-color: #f8f8f8;
  }

  .editor {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    padding: 8px;
    border: 2px solid #1976d2;
    box-sizing: border-box;
    outline: none;
  }

  .selected-content {
    padding: 12px;
    background: #f5f5f5;
    border-top: 1px solid #ddd;
    font-size: 0.9rem;
  }
</style>
