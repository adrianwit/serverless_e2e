import React from 'react'
import { Table } from 'semantic-ui-react'
import PropTypes from 'prop-types'

const ProductsList = ({ title, children }) => (
    <div id="productList">
    <h3>{title}</h3>
    <Table striped>
        <Table.Header>
        <Table.Row>
        <Table.HeaderCell>Name</Table.HeaderCell>
        <Table.HeaderCell textAlign='right'>Price</Table.HeaderCell>
        <Table.HeaderCell textAlign='right'>Quantity</Table.HeaderCell>
        <Table.HeaderCell> </Table.HeaderCell>
        </Table.Row>
        </Table.Header>
        {children}
    </Table>

  </div>
)

ProductsList.propTypes = {
  children: PropTypes.node,
  title: PropTypes.string.isRequired
}

export default ProductsList
