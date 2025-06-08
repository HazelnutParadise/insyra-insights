<script lang="ts">
  import { onMount, createEventDispatcher } from "svelte";
  import { TableData, EditingState } from "../types/datatable";
  import {
    GetTableData,
    UpdateCellValue,
    UpdateColumnName,
  } from "../../wailsjs/go/main/App"; // 組件屬性
  export let tableName: string;
  export let tableKey: number = 0; // 用於強制重新載入的 key

  // 創建事件分發器
  const dispatch = createEventDispatcher();

  // 狀態變數
  let tableData: TableData | null = null;
  let loading = true;
  let error = "";
  // 編輯狀態
  let editingState: EditingState = {
    tableName: "",
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
  // 編輯輸入元素引用
  let editInput: HTMLInputElement;

  // 當進入編輯模式時，設置焦點
  $: if (editingState.isEditing && editInput) {
    setTimeout(() => {
      editInput.focus();
      editInput.select();
    }, 0);
  }

  onMount(async () => {
    await loadTableData();
  }); // 當 tableName 或 tableKey 變化時重新載入
  $: if (tableName || tableKey >= 0) {
    loadTableData();
  }

  // 載入表格資料
  async function loadTableData() {
    try {
      loading = true;
      const data = await GetTableData(tableName);
      tableData = data as TableData;

      // 計算並分發統計數據
      if (tableData) {
        const stats = calculateStatistics(tableData);
        dispatch("statsUpdate", stats);
      }
    } catch (err) {
      error = `載入資料表失敗: ${err}`;
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
  }
  // 儲存格點擊處理
  function handleCellClick(rowIndex: number, colIndex: number, value: string) {
    selectedRow = rowIndex;
    selectedCol = colIndex;
    selectedCellContent = value;

    // 如果正在編輯其他儲存格，先結束編輯
    if (editingState.isEditing) {
      handleEditComplete();
    }
  }

  // 儲存格雙擊處理 (進入編輯模式)
  function handleCellDblClick(
    rowIndex: number,
    colIndex: number,
    colName: string,
    value: string
  ) {
    editingState = {
      tableName,
      rowIndex,
      colIndex,
      colName,
      value,
      isEditing: true,
    };
  }

  // 欄位標題點擊處理
  function handleColumnHeaderClick(colIndex: number, colName: string) {
    selectedCol = colIndex;
    selectedRow = -1;
    selectedCellContent = colName;

    // 如果正在編輯其他儲存格，先結束編輯
    if (editingState.isEditing) {
      handleEditComplete();
    }
  }
  // 欄位標題雙擊處理 (進入編輯模式)
  function handleColumnHeaderDblClick(colIndex: number, colName: string) {
    editingState = {
      tableName,
      rowIndex: -1,
      colIndex,
      colName,
      value: colName,
      isEditing: true,
    };
  }

  // 編輯完成處理
  async function handleEditComplete() {
    if (!editingState.isEditing) return;
    try {
      if (editingState.rowIndex >= 0) {
        // 更新儲存格值
        await UpdateCellValue(
          tableName,
          editingState.rowIndex,
          editingState.colIndex,
          editingState.value
        );
      } else {
        // 更新欄名
        await UpdateColumnName(
          tableName,
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
        tableName: "",
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
          <tr>
            <!-- 空白頂角儲存格 -->
            <th class="corner-cell"></th>

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
            <tr class:selected-row={rowIndex === selectedRow}>
              <!-- 行標識 -->
              <td class="row-header">{rowIndex + 1}</td>

              <!-- 儲存格資料 -->
              {#each tableData.columns as column, colIndex}
                {@const cellValue = row.cells[column.name] || ""}
                <td
                  class="cell"
                  class:selected-cell={rowIndex === selectedRow &&
                    colIndex === selectedCol}
                  class:selected-col={colIndex === selectedCol}
                  on:click={() =>
                    handleCellClick(rowIndex, colIndex, String(cellValue))}
                  on:dblclick={() =>
                    handleCellDblClick(
                      rowIndex,
                      colIndex,
                      column.name,
                      String(cellValue)
                    )}
                  on:keydown={(e) => {
                    if (e.key === "Enter" || e.key === " ") {
                      handleCellClick(rowIndex, colIndex, String(cellValue));
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
                    {cellValue}
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
    width: 100%;
    border-collapse: collapse;
    table-layout: fixed;
  }

  th,
  td {
    padding: 8px 12px;
    border: 1px solid #ddd;
    text-align: left;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    min-width: 100px;
  }

  .corner-cell {
    background-color: #f0f0f0;
    width: 50px;
    min-width: 50px;
  }

  .column-header {
    background-color: #e1eeff;
    position: sticky;
    top: 0;
    z-index: 10;
    font-weight: 600;
    transition: background-color 0.2s;
  }

  .column-header.selected {
    background-color: #c8e6c9;
  }

  .row-header {
    background-color: #f0f0f0;
    position: sticky;
    left: 0;
    z-index: 5;
    width: 50px;
    min-width: 50px;
  }

  .selected-row {
    background-color: rgba(200, 230, 201, 0.3);
  }

  .selected-col {
    background-color: rgba(200, 230, 201, 0.3);
  }

  .selected-cell {
    background-color: rgba(200, 230, 201, 0.6);
  }

  .cell {
    position: relative;
    transition: background-color 0.2s;
  }

  .cell:hover {
    background-color: #f9f9f9;
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
